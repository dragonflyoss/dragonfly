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

package service

import (
	"context"
	"fmt"

	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/types"
	"d7y.io/dragonfly/v2/pkg/structure"
)

// convertToMap converts a struct to map, returns nil if the input is nil.
func convertToMap(v any) (map[string]any, error) {
	if v == nil {
		return nil, nil
	}
	return structure.StructToMap(v)
}

func (s *service) CreateSchedulerCluster(ctx context.Context, json types.CreateSchedulerClusterRequest) (*models.SchedulerCluster, error) {
	config, err := convertToMap(json.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to convert config: %w", err)
	}

	clientConfig, err := convertToMap(json.ClientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to convert client config: %w", err)
	}

	scopes, err := convertToMap(json.Scopes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert scopes: %w", err)
	}

	schedulerCluster := models.SchedulerCluster{
		Name:         json.Name,
		BIO:          json.BIO,
		Config:       config,
		ClientConfig: clientConfig,
		Scopes:       scopes,
		IsDefault:    json.IsDefault,
	}

	if err := s.db.WithContext(ctx).Create(&schedulerCluster).Error; err != nil {
		return nil, fmt.Errorf("failed to create scheduler cluster: %w", err)
	}

	if json.SeedPeerClusterID > 0 {
		if err := s.AddSchedulerClusterToSeedPeerCluster(ctx, json.SeedPeerClusterID, schedulerCluster.ID); err != nil {
			return nil, fmt.Errorf("failed to add scheduler cluster to seed peer cluster: %w", err)
		}
	}

	return &schedulerCluster, nil
}

func (s *service) DestroySchedulerCluster(ctx context.Context, id uint) error {
	schedulerCluster := models.SchedulerCluster{}
	if err := s.db.WithContext(ctx).Preload("Schedulers").First(&schedulerCluster, id).Error; err != nil {
		return fmt.Errorf("failed to find scheduler cluster with id %d: %w", id, err)
	}

	if len(schedulerCluster.Schedulers) != 0 {
		return fmt.Errorf("scheduler cluster %d cannot be deleted: it still has %d scheduler(s)", id, len(schedulerCluster.Schedulers))
	}

	if err := s.db.WithContext(ctx).Unscoped().Delete(&models.SchedulerCluster{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete scheduler cluster %d: %w", id, err)
	}

	return nil
}

func (s *service) UpdateSchedulerCluster(ctx context.Context, id uint, json types.UpdateSchedulerClusterRequest) (*models.SchedulerCluster, error) {
	config, err := convertToMap(json.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to convert config: %w", err)
	}

	clientConfig, err := convertToMap(json.ClientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to convert client config: %w", err)
	}

	scopes, err := convertToMap(json.Scopes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert scopes: %w", err)
	}

	schedulerCluster := models.SchedulerCluster{}
	if err := s.db.WithContext(ctx).First(&schedulerCluster, id).Error; err != nil {
		return nil, fmt.Errorf("failed to find scheduler cluster with id %d: %w", id, err)
	}

	// Build update map, only include non-nil fields
	// Note: For string fields, empty string means "don't update" due to binding:"omitempty"
	updates := make(map[string]any)
	if json.Name != "" {
		updates["name"] = json.Name
	}
	if json.BIO != "" {
		updates["bio"] = json.BIO
	}
	if config != nil {
		updates["config"] = config
	}
	if clientConfig != nil {
		updates["client_config"] = clientConfig
	}
	if scopes != nil {
		updates["scopes"] = scopes
	}

	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(&schedulerCluster).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update scheduler cluster %d: %w", id, err)
		}
	}

	// Updates does not accept bool as false.
	// Refer to https://stackoverflow.com/questions/56653423/gorm-doesnt-update-boolean-field-to-false.
	if json.IsDefault != schedulerCluster.IsDefault {
		if err := s.db.WithContext(ctx).Model(&schedulerCluster).Update("is_default", json.IsDefault).Error; err != nil {
			return nil, fmt.Errorf("failed to update is_default for scheduler cluster %d: %w", id, err)
		}
	}

	if json.SeedPeerClusterID > 0 {
		if err := s.AddSchedulerClusterToSeedPeerCluster(ctx, json.SeedPeerClusterID, schedulerCluster.ID); err != nil {
			return nil, fmt.Errorf("failed to add scheduler cluster to seed peer cluster: %w", err)
		}
	}

	// Reload to get updated data
	if err := s.db.WithContext(ctx).Preload("SeedPeerClusters").First(&schedulerCluster, id).Error; err != nil {
		return nil, fmt.Errorf("failed to reload scheduler cluster %d: %w", id, err)
	}

	return &schedulerCluster, nil
}

func (s *service) GetSchedulerCluster(ctx context.Context, id uint) (*models.SchedulerCluster, error) {
	schedulerCluster := models.SchedulerCluster{}
	if err := s.db.WithContext(ctx).Preload("SeedPeerClusters").First(&schedulerCluster, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get scheduler cluster with id %d: %w", id, err)
	}

	return &schedulerCluster, nil
}

func (s *service) GetSchedulerClusters(ctx context.Context, q types.GetSchedulerClustersQuery) ([]models.SchedulerCluster, int64, error) {
	var count int64
	var schedulerClusters []models.SchedulerCluster

	// Build query
	query := s.db.WithContext(ctx)
	if q.Name != "" {
		query = query.Where("name = ?", q.Name)
	}

	// Count total records (before pagination)
	if err := query.Model(&models.SchedulerCluster{}).Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count scheduler clusters: %w", err)
	}

	// Apply pagination and preload
	if err := query.Scopes(models.Paginate(q.Page, q.PerPage)).
		Preload("SeedPeerClusters").
		Find(&schedulerClusters).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get scheduler clusters: %w", err)
	}

	return schedulerClusters, count, nil
}

func (s *service) AddSchedulerToSchedulerCluster(ctx context.Context, id, schedulerID uint) error {
	schedulerCluster := models.SchedulerCluster{}
	if err := s.db.WithContext(ctx).First(&schedulerCluster, id).Error; err != nil {
		return fmt.Errorf("failed to find scheduler cluster with id %d: %w", id, err)
	}

	scheduler := models.Scheduler{}
	if err := s.db.WithContext(ctx).First(&scheduler, schedulerID).Error; err != nil {
		return fmt.Errorf("failed to find scheduler with id %d: %w", schedulerID, err)
	}

	if err := s.db.WithContext(ctx).Model(&schedulerCluster).Association("Schedulers").Append(&scheduler); err != nil {
		return fmt.Errorf("failed to add scheduler %d to cluster %d: %w", schedulerID, id, err)
	}

	return nil
}
