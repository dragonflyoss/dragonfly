/*
 *     Copyright 2025 The Dragonfly Authors
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

package gc

import (
	"errors"

	"gorm.io/gorm"

	"d7y.io/dragonfly/v2/manager/models"
)

const (
	// GCJobType indicates the gc task is completed successfully.
	GCJobType = "gc"

	// GCStateSuccess indicates the gc task is completed successfully.
	GCStateSuccess = "SUCCESS"

	// GCStateFailure indicates the gc task is completed with failure.
	GCStateFailure = "FAILURE"
)

type Result struct {
	Error  error
	Purged int64
}

type jobRecorder struct {
	db  *gorm.DB
	job *models.Job
}

func newJobRecorder(db *gorm.DB) *jobRecorder {
	return &jobRecorder{
		db: db,
	}
}

func (jb *jobRecorder) Init(userID uint, taskID string, args models.JSONMap) error {
	job := models.Job{
		Type:   GCJobType,
		TaskID: taskID,
		UserID: userID,
		Args:   args,
	}

	if err := jb.db.Create(&job).Error; err != nil {
		return err
	}

	jb.job = &job
	return nil
}

func (jb *jobRecorder) Record(result Result) error {
	if jb.job == nil {
		return errors.New("job not found")
	}

	if jb.job.Result == nil {
		jb.job.Result = make(models.JSONMap)
	}

	jb.job.State = GCStateSuccess
	jb.job.Result["purged"] = result.Purged

	if result.Error != nil {
		jb.job.State = GCStateFailure
		jb.job.Result["error"] = result.Error.Error()
	}

	return jb.db.Save(jb.job).Error
}
