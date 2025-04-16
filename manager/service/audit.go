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

package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/types"
)

const (
	// auditBatchInsertSize is the size for batch insertion.
	auditBatchInsertSize = 500
	// auditBatchInsertInterval is the interval for batch insertion.
	auditBatchInsertInterval = 5 * time.Second
	// auditBufferSize is the size of the audit channel buffer.
	auditBufferSize = 1000
)

var (
	auditCh chan *models.Audit
	once    sync.Once
)

func (s *service) AsyncCreateAudit(ctx context.Context, json *types.CreateAuditRequest) error {
	once.Do(func() {
		auditCh = make(chan *models.Audit, auditBufferSize)
		go s.processAudit()
	})

	audit := &models.Audit{
		ActorType:  json.ActorType,
		ActorName:  json.ActorName,
		EventType:  json.EventType,
		Operation:  json.Operation,
		OperatedAt: json.OperatedAt,
		State:      json.State,
		Path:       json.Path,
		StatusCode: json.StatusCode,
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("context done: %w", ctx.Err())
	default:
		select {
		case auditCh <- audit:
			return nil
		default:
			// Avoid to hang out the AsyncCreateAudit if the buffer is full.
			return fmt.Errorf("audit buffer is full, buffer size: %d, drop the audit %#v", auditBufferSize, audit)
		}
	}
}

func (s *service) processAudit() {
	// Use the new context as this is asynchronous operation.
	ctx := context.Background()
	audits := make([]*models.Audit, 0, auditBatchInsertSize)
	ticker := time.NewTicker(auditBatchInsertInterval)
	defer ticker.Stop()

	createAuditInBatch := func(ctx context.Context, audits []*models.Audit) error {
		if err := s.db.WithContext(ctx).CreateInBatches(&audits, len(audits)).Error; err != nil {
			logger.Errorf("failed to create audit in batch: %v", err)
			return err
		}

		return nil
	}

	for {
		select {
		case audit, ok := <-auditCh:
			if !ok {
				// Channel closed.
				if len(audits) > 0 {
					createAuditInBatch(ctx, audits) //nolint
					return
				}
				return
			}

			audits = append(audits, audit)
			if len(audits) >= auditBatchInsertSize {
				if err := createAuditInBatch(ctx, audits); err == nil {
					audits = make([]*models.Audit, 0, auditBatchInsertSize)
				}

				ticker.Reset(auditBatchInsertInterval)
			}
		case <-ticker.C:
			if len(audits) > 0 {
				if err := createAuditInBatch(ctx, audits); err == nil {
					audits = make([]*models.Audit, 0, auditBatchInsertSize)
				}
			}
		}
	}
}

func (s *service) GetAudits(ctx context.Context, q types.GetAuditsQuery) ([]models.Audit, int64, error) {
	var count int64
	audits := []models.Audit{}
	if err := s.db.WithContext(ctx).Scopes(models.Paginate(q.Page, q.PerPage)).Find(&audits).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return audits, count, nil
}
