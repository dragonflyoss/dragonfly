package localcdn

import (
	"context"
	"github.com/dragonflyoss/Dragonfly2/cdnsystem/source"
	"github.com/dragonflyoss/Dragonfly2/cdnsystem/types"
	"path"

	"github.com/dragonflyoss/Dragonfly2/cdnsystem/config"
	"github.com/dragonflyoss/Dragonfly2/cdnsystem/daemon/mgr"
	"github.com/dragonflyoss/Dragonfly2/cdnsystem/store"
	"github.com/dragonflyoss/Dragonfly2/cdnsystem/util"
	"github.com/dragonflyoss/Dragonfly2/pkg/limitreader"
	"github.com/dragonflyoss/Dragonfly2/pkg/metricsutils"
	"github.com/dragonflyoss/Dragonfly2/pkg/ratelimiter"
	"github.com/dragonflyoss/Dragonfly2/pkg/stringutils"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var _ mgr.CDNMgr = &Manager{}

func init() {
	mgr.Register(config.CDNPatternLocal, NewManager)
}

type metrics struct {
	cdnCacheHitCount     *prometheus.CounterVec
	cdnDownloadCount     *prometheus.CounterVec
	cdnDownloadBytes     *prometheus.CounterVec
	cdnDownloadFailCount *prometheus.CounterVec
}

func newMetrics(register prometheus.Registerer) *metrics {
	return &metrics{
		cdnCacheHitCount: metricsutils.NewCounter(config.SubsystemSupernode, "cdn_cache_hit_total",
			"Total times of hitting cdn cache", []string{}, register),

		cdnDownloadCount: metricsutils.NewCounter(config.SubsystemSupernode, "cdn_download_total",
			"Total times of cdn download", []string{}, register),

		cdnDownloadBytes: metricsutils.NewCounter(config.SubsystemSupernode, "cdn_download_size_bytes_total",
			"total file size of cdn downloaded from source in bytes", []string{}, register,
		),

		cdnDownloadFailCount: metricsutils.NewCounter(config.SubsystemSupernode, "cdn_download_failed_total",
			"Total failure times of cdn download", []string{}, register),
	}
}

// Manager is an implementation of the interface of CDNMgr.
type Manager struct {
	cfg                  *config.Config
	cacheStore           *store.Store
	limiter              *ratelimiter.RateLimiter
	cdnLocker            *util.LockerPool
	metaDataManager      *fileMetaDataManager
	pieceMetaDataManager *pieceMetaDataManager
	cdnReporter          *reporter
	detector             *cacheDetector
	resourceClient       source.ResourceClient
	writer               *cacheWriter
	metrics              *metrics
}

// NewManager returns a new Manager.
func NewManager(cfg *config.Config, cacheStore *store.Store, resourceClient source.ResourceClient,
	rateLimiter *ratelimiter.RateLimiter, register prometheus.Registerer) (mgr.CDNMgr, error) {
	return newManager(cfg, cacheStore, resourceClient, rateLimiter, register)
}

func newManager(cfg *config.Config, cacheStore *store.Store,
	resourceClient source.ResourceClient, rateLimiter *ratelimiter.RateLimiter, register prometheus.Registerer) (*Manager, error) {
	metaDataManager := newFileMetaDataManager(cacheStore)
	pieceMetaDataManager := newPieceMetaDataMgr()
	cdnReporter := newReporter(pieceMetaDataManager)
	return &Manager{
		cfg:                  cfg,
		cacheStore:           cacheStore,
		limiter:              rateLimiter,
		cdnLocker:            util.NewLockerPool(),
		metaDataManager:      metaDataManager,
		pieceMetaDataManager: pieceMetaDataManager,
		cdnReporter:          cdnReporter,
		detector:             newCacheDetector(cacheStore, metaDataManager, pieceMetaDataManager, resourceClient),
		resourceClient:       resourceClient,
		writer:               newCacheWriter(cacheStore, cdnReporter),
		metrics:              newMetrics(register),
	}, nil
}

// TriggerCDN will trigger CDN to download the file from sourceUrl.
func (cm *Manager) TriggerCDN(ctx context.Context, task *types.SeedTaskInfo) (*types.SeedTaskInfo, error) {
	sourceFileLength := task.SourceFileLength
	if sourceFileLength == 0 {
		sourceFileLength = -1
	}
	// obtain taskID write lock
	cm.cdnLocker.GetLock(task.TaskID, false)
	defer cm.cdnLocker.ReleaseLock(task.TaskID, false)
	// detect Cache
	detectResult, err := cm.detector.detectCache(ctx, task)
	if err != nil {
		logrus.Errorf("taskId: %s failed to detect cache err: %v", task.TaskID, err)
	}
	// report cache
	updateTaskInfo, err := cm.cdnReporter.reportCache(task.TaskID, detectResult)
	if err != nil {
		logrus.Errorf("taskId: %s failed to report cache err: %v", task.TaskID, err)
	}

	if detectResult.breakNum == -1 {
		logrus.Infof("taskId: %s cache full hit on local", task.TaskID)
		cm.metrics.cdnCacheHitCount.WithLabelValues().Inc()
		return updateTaskInfo, nil
	}

	// start to download the source file
	resp, err := cm.download(ctx, task.TaskID, task.Url, task.Headers, detectResult.breakNum, sourceFileLength, task.PieceSize)
	cm.metrics.cdnDownloadCount.WithLabelValues().Inc()
	if err != nil {
		cm.metrics.cdnDownloadFailCount.WithLabelValues().Inc()
		return getUpdateTaskInfoWithStatusOnly(types.TaskInfoCdnStatusFAILED), err
	}
	defer resp.Body.Close()

	cm.updateExpireInfo(ctx, task.TaskID, resp.ExpireInfo)
	reader := limitreader.NewLimitReaderWithLimiterAndMD5Sum(resp.Body, cm.limiter, detectResult.fileMd5)
	downloadMetadata, err := cm.writer.startWriter(ctx, reader, task, detectResult.breakNum, detectResult.downloadedFileLength, sourceFileLength)
	if err != nil {
		logrus.Errorf("failed to write for task %s: %v", task.TaskID, err)
		return getUpdateTaskInfoWithStatusOnly(types.TaskInfoCdnStatusFAILED), err
	}
	cm.metrics.cdnDownloadBytes.WithLabelValues().Add(float64(downloadMetadata.realSourceFileLength))

	realMD5 := reader.Md5()
	success, err := cm.handleCDNResult(ctx, task, realMD5, sourceFileLength, downloadMetadata.realSourceFileLength, downloadMetadata.realCdnFileLength)
	if err != nil || !success {
		return getUpdateTaskInfoWithStatusOnly(types.TaskInfoCdnStatusFAILED), err
	}

	return getUpdateTaskInfo(types.TaskInfoCdnStatusSUCCESS, realMD5, downloadMetadata.realCdnFileLength), nil
}

// GetHTTPPath returns the http download path of taskID.
// The returned path joined the DownloadRaw.Bucket and DownloadRaw.Key.
func (cm *Manager) GetHTTPPath(ctx context.Context, task *types.SeedTaskInfo) (string, error) {
	raw := getDownloadRawFunc(task.TaskID)
	return path.Join("/", raw.Bucket, raw.Key), nil
}

// GetStatus gets the status of the file.
func (cm *Manager) GetStatus(ctx context.Context, taskID string) (cdnStatus string, err error) {
	return types.TaskInfoCdnStatusSUCCESS, nil
}

// CheckFile checks the file whether exists.
func (cm *Manager) CheckFileExist(ctx context.Context, taskID string) bool {
	if _, err := cm.cacheStore.Stat(ctx, getDownloadRaw(taskID)); err != nil {
		return false
	}
	return true
}

// Delete the cdn meta with specified taskID.
// It will also delete the files on the disk when the force equals true.
func (cm *Manager) Delete(ctx context.Context, taskID string, force bool) error {
	if !force {
		return cm.pieceMetaDataManager.removePieceMetaRecordsByTaskID(taskID)
	}

	return deleteTaskFiles(ctx, cm.cacheStore, taskID)
}

func (cm *Manager) handleCDNResult(ctx context.Context, task *types.SeedTaskInfo, realMd5 string, sourceFileLength, realSourceFileLength, realCdnFileLength int64) (bool, error) {
	var isSuccess = true
	if !stringutils.IsEmptyStr(task.RequestMd5) && task.RequestMd5 != realMd5 {
		logrus.Errorf("taskId:%s url:%s file md5 not match expected:%s real:%s", task.TaskID, task.Url, task.RequestMd5, realMd5)
		isSuccess = false
	}
	if isSuccess && sourceFileLength >= 0 && sourceFileLength != realSourceFileLength {
		logrus.Errorf("taskId:%s url:%s file length not match expected:%d real:%d", task.TaskID, task.Url, sourceFileLength, realSourceFileLength)
		isSuccess = false
	}

	if !isSuccess {
		realCdnFileLength = 0
	}
	if err := cm.metaDataManager.updateStatusAndResult(ctx, task.TaskID, &fileMetaData{
		Finish:        true,
		Success:       isSuccess,
		Md5:           realMd5,
		CdnFileLength: realCdnFileLength,
	}); err != nil {
		return false, err
	}

	if !isSuccess {
		return false, nil
	}

	logrus.Infof("success to get taskID: %s fileLength: %d realMd5: %s", task.TaskID, realCdnFileLength, realMd5)

	pieceMD5s, err := cm.pieceMetaDataManager.getPieceMetaRecordsByTaskID(task.TaskID)
	if err != nil {
		return false, err
	}

	if err := cm.pieceMetaDataManager.writePieceMetaRecords(ctx, task.TaskID, realMd5, pieceMD5s); err != nil {
		return false, err
	}
	return true, nil
}

func (cm *Manager) updateExpireInfo(ctx context.Context, taskID string, expireInfo map[string]string) {
	if err := cm.metaDataManager.updateExpireInfo(ctx, taskID, expireInfo); err != nil {
		logrus.Errorf("taskID: %s failed to update expireInfo(%s): %v", taskID, expireInfo, err)
	}
	logrus.Infof("taskID: %s success to update expireInfo(%s)", expireInfo, taskID)
}
