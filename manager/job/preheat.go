/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

//go:generate mockgen -destination mocks/preheat_mock.go -source preheat.go -package mocks

package job

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"time"

	machineryv1tasks "github.com/dragonflyoss/machinery/v1/tasks"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	internaljob "d7y.io/dragonfly/v2/internal/job"
	"d7y.io/dragonfly/v2/manager/config"
	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/types"
)

// preheatImage is an image for preheat.
type PreheatType string

// String returns the string representation of PreheatType.
func (p PreheatType) String() string {
	return string(p)
}

const (
	// PreheatImageType is image type of preheat job.
	PreheatImageType PreheatType = "image"

	// PreheatFileType is file type of preheat job.
	PreheatFileType PreheatType = "file"
)

// Preheat is an interface for preheat job.
type Preheat interface {
	// CreatePreheat creates a preheat job.
	CreatePreheat(context.Context, []models.Scheduler, types.PreheatArgs) (*internaljob.GroupJobState, error)
}

// preheat is an implementation of Preheat.
type preheat struct {
	job                *internaljob.Job
	internalJobImage   internaljob.Image
	rootCAs            *x509.CertPool
	insecureSkipVerify bool
}

// newPreheat creates a new Preheat.
func newPreheat(job *internaljob.Job, internalJobImage internaljob.Image, rootCAs *x509.CertPool, insecureSkipVerify bool) Preheat {
	return &preheat{
		job:                job,
		internalJobImage:   internalJobImage,
		rootCAs:            rootCAs,
		insecureSkipVerify: insecureSkipVerify,
	}
}

// CreatePreheat creates a preheat job.
func (p *preheat) CreatePreheat(ctx context.Context, schedulers []models.Scheduler, json types.PreheatArgs) (*internaljob.GroupJobState, error) {
	var span trace.Span
	ctx, span = tracer.Start(ctx, config.SpanPreheat, trace.WithSpanKind(trace.SpanKindProducer))
	span.SetAttributes(config.AttributePreheatType.String(json.Type))
	span.SetAttributes(config.AttributePreheatURL.String(json.URL))
	defer span.End()

	// Generate download files.
	var files []*internaljob.PreheatRequest
	var err error
	switch PreheatType(json.Type) {
	case PreheatImageType:
		// Image preheat only supports to preheat single image.
		if json.URL == "" {
			return nil, errors.New("invalid params: url is required")
		}

		files, err = p.internalJobImage.CreatePreheatRequestsByManifestURL(ctx, &internaljob.ManifestRequest{
			URL:                         json.URL,
			PieceLength:                 json.PieceLength,
			Tag:                         json.Tag,
			Application:                 json.Application,
			FilteredQueryParams:         json.FilteredQueryParams,
			Headers:                     json.Headers,
			Username:                    json.Username,
			Password:                    json.Password,
			Platform:                    json.Platform,
			Scope:                       json.Scope,
			IPs:                         json.IPs,
			Percentage:                  json.Percentage,
			Count:                       json.Count,
			ConcurrentTaskCount:         json.ConcurrentTaskCount,
			ConcurrentPeerCount:         json.ConcurrentPeerCount,
			Timeout:                     json.Timeout,
			RootCAs:                     p.rootCAs,
			InsecureSkipVerify:          p.insecureSkipVerify,
			EnableTaskIDBasedBlobDigest: json.EnableTaskIDBasedBlobDigest,
		})
		if err != nil {
			return nil, err
		}
	case PreheatFileType:
		urls := json.URLs
		if json.URL != "" {
			urls = append(urls, json.URL)
		}

		if len(urls) == 0 {
			return nil, errors.New("invalid params: url or urls is required")
		}

		var certificateChain [][]byte
		if p.rootCAs != nil {
			certificateChain = p.rootCAs.Subjects()
		}

		files = append(files, &internaljob.PreheatRequest{
			URLs:                        urls,
			PieceLength:                 json.PieceLength,
			Tag:                         json.Tag,
			Application:                 json.Application,
			FilteredQueryParams:         json.FilteredQueryParams,
			Headers:                     json.Headers,
			Scope:                       json.Scope,
			IPs:                         json.IPs,
			Percentage:                  json.Percentage,
			Count:                       json.Count,
			ConcurrentTaskCount:         json.ConcurrentTaskCount,
			ConcurrentPeerCount:         json.ConcurrentPeerCount,
			CertificateChain:            certificateChain,
			InsecureSkipVerify:          p.insecureSkipVerify,
			Timeout:                     json.Timeout,
			EnableTaskIDBasedBlobDigest: json.EnableTaskIDBasedBlobDigest,
		})

	default:
		return nil, errors.New("unknown preheat type")
	}

	// Initialize queues.
	queues, err := getSchedulerQueues(schedulers)
	if err != nil {
		return nil, err
	}

	return p.createGroupJob(ctx, files, queues)
}

// createGroupJob creates a group job.
func (p *preheat) createGroupJob(ctx context.Context, files []*internaljob.PreheatRequest, queues []internaljob.Queue) (*internaljob.GroupJobState, error) {
	groupUUID := fmt.Sprintf("group_%s", uuid.New().String())
	var signatures []*machineryv1tasks.Signature
	for _, queue := range queues {
		for _, file := range files {
			file.GroupUUID = groupUUID
			taskUUID := fmt.Sprintf("task_%s", uuid.New().String())
			file.TaskUUID = taskUUID

			args, err := internaljob.MarshalRequest(file)
			if err != nil {
				logger.Errorf("[preheat]: preheat marshal request: %v, error: %v", file, err)
				continue
			}

			signatures = append(signatures, &machineryv1tasks.Signature{
				UUID:       taskUUID,
				Name:       internaljob.PreheatJob,
				RoutingKey: queue.String(),
				Args:       args,
			})
		}
	}

	group, err := machineryv1tasks.NewGroup(signatures...)
	if err != nil {
		return nil, err
	}
	group.GroupUUID = groupUUID

	var tasks []machineryv1tasks.Signature
	for _, signature := range signatures {
		tasks = append(tasks, *signature)
	}

	logger.Infof("[preheat]: create preheat group %s in queues %v, tasks: %#v", group.GroupUUID, queues, tasks)
	if _, err := p.job.Server.SendGroupWithContext(ctx, group, 50); err != nil {
		logger.Errorf("[preheat]: create preheat group %s failed", group.GroupUUID, err)
		return nil, err
	}

	return &internaljob.GroupJobState{
		GroupUUID: group.GroupUUID,
		State:     machineryv1tasks.StatePending,
		CreatedAt: time.Now(),
	}, nil
}
