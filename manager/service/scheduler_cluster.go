package service

import (
	"d7y.io/dragonfly/v2/manager/model"
	"d7y.io/dragonfly/v2/manager/types"
)

func (s *service) CreateSchedulerCluster(json types.CreateSchedulerClusterRequest) (*model.SchedulerCluster, error) {
	schedulerCluster := model.SchedulerCluster{
		Name:         json.Name,
		BIO:          json.BIO,
		Config:       json.Config,
		ClientConfig: json.ClientConfig,
	}

	if err := s.db.Create(&schedulerCluster).Error; err != nil {
		return nil, err
	}

	return &schedulerCluster, nil
}

func (s *service) DestroySchedulerCluster(id uint) error {
	if err := s.db.Unscoped().Delete(&model.SchedulerCluster{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateSchedulerCluster(id uint, json types.UpdateSchedulerClusterRequest) (*model.SchedulerCluster, error) {
	schedulerCluster := model.SchedulerCluster{}
	if err := s.db.First(&schedulerCluster, id).Updates(model.SchedulerCluster{
		Name:         json.Name,
		BIO:          json.BIO,
		Config:       json.Config,
		ClientConfig: json.ClientConfig,
	}).Error; err != nil {
		return nil, err
	}

	return &schedulerCluster, nil
}

func (s *service) GetSchedulerCluster(id uint) (*model.SchedulerCluster, error) {
	schedulerCluster := model.SchedulerCluster{}
	if err := s.db.First(&schedulerCluster, id).Error; err != nil {
		return nil, err
	}

	return &schedulerCluster, nil
}

func (s *service) GetSchedulerClusters(q types.GetSchedulerClustersQuery) (*[]model.SchedulerCluster, error) {
	schedulerClusters := []model.SchedulerCluster{}
	if err := s.db.Scopes(model.Paginate(q.Page, q.PerPage)).Where(&model.SchedulerCluster{
		Name: q.Name,
	}).Find(&schedulerClusters).Error; err != nil {
		return nil, err
	}

	return &schedulerClusters, nil
}

func (s *service) SchedulerClusterTotalCount(q types.GetSchedulerClustersQuery) (int64, error) {
	var count int64
	if err := s.db.Model(&model.SchedulerCluster{}).Where(&model.SchedulerCluster{
		Name: q.Name,
	}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *service) AddSchedulerToSchedulerCluster(id, schedulerID uint) error {
	schedulerCluster := model.SchedulerCluster{}
	if err := s.db.First(&schedulerCluster, id).Error; err != nil {
		return err
	}

	scheduler := model.Scheduler{}
	if err := s.db.First(&scheduler, schedulerID).Error; err != nil {
		return err
	}

	if err := s.db.Model(&schedulerCluster).Association("Schedulers").Append(&scheduler); err != nil {
		return err
	}

	return nil
}
