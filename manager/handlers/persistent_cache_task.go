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

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"d7y.io/dragonfly/v2/manager/types"
)

// @Summary Destroy PersistentCacheTask
// @Description Destroy PersistentCacheTask by id
// @Tags PersistentCacheTask
// @Accept json
// @Produce json
// @Param scheduler_cluster_id query uint true "scheduler cluster id"
// @Param id path string true "task id"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /api/v1/persistent-cache-tasks/{id} [delete]
func (h *Handlers) DestroyPersistentCacheTask(ctx *gin.Context) {
	var params types.PersistentCacheTaskParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	var query types.DestroyPersistentCacheTaskQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	if err := h.service.DestroyPersistentCacheTask(ctx.Request.Context(), query.SchedulerClusterID, params.ID); err != nil {
		ctx.Error(err) // nolint: errcheck
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Get PersistentCacheTask
// @Description Get PersistentCacheTask by id
// @Tags PersistentCacheTask
// @Accept json
// @Produce json
// @Param scheduler_cluster_id query uint true "scheduler cluster id"
// @Param id path string true "task id"
// @Success 200 {object} types.PersistentCacheTask
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /api/v1/persistent-cache-tasks/{id} [get]
func (h *Handlers) GetPersistentCacheTask(ctx *gin.Context) {
	var params types.PersistentCacheTaskParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	var query types.GetPersistentCacheTaskQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	persistentCacheTask, err := h.service.GetPersistentCacheTask(ctx.Request.Context(), query.SchedulerClusterID, params.ID)
	if err != nil {
		ctx.Error(err) // nolint: errcheck
		return
	}

	ctx.JSON(http.StatusOK, persistentCacheTask)
}

// @Summary Get PersistentCacheTasks
// @Description Get PersistentCacheTasks
// @Tags PersistentCacheTask
// @Accept json
// @Produce json
// @Param scheduler_cluster_id query uint true "scheduler cluster id"
// @Param page query int true "current page" default(0)
// @Param per_page query int true "return max item count, default 10, max 50" default(10) minimum(2) maximum(50)
// @Success 200 {object} []types.PersistentCacheTask
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /api/v1/persistent-cache-tasks [get]
func (h *Handlers) GetPersistentCacheTasks(ctx *gin.Context) {
	var query types.GetPersistentCacheTasksQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	h.setPaginationDefault(&query.Page, &query.PerPage)
	persistentCacheTasks, count, err := h.service.GetPersistentCacheTasks(ctx.Request.Context(), query)
	if err != nil {
		ctx.Error(err) // nolint: errcheck
		return
	}

	h.setPaginationLinkHeader(ctx, query.Page, query.PerPage, int(count))
	ctx.JSON(http.StatusOK, persistentCacheTasks)
}
