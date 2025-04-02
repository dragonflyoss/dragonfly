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

// @Summary Destroy PersistentCache
// @Description Destroy PersistentCache by id
// @Tags PersistentCache
// @Accept json
// @Produce json
// @Param scheduler_cluster_id path string true "scheduler cluster id"
// @Param task_id path string true "task id"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /api/v1/persistent-caches/{scheduler_cluster_id}/{task_id} [delete]
func (h *Handlers) DestroyPersistentCache(ctx *gin.Context) {
	var params types.PersistentCacheParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	if err := h.service.DestroyPersistentCache(ctx.Request.Context(), params); err != nil {
		ctx.Error(err) // nolint: errcheck
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Get PersistentCache
// @Description Get PersistentCache by id
// @Tags PersistentCache
// @Accept json
// @Produce json
// @Param scheduler_cluster_id path string true "scheduler cluster id"
// @Param task_id path string true "task id"
// @Success 200 {object} types.GetPersistentCacheResponse
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /api/v1/persistent-caches/{scheduler_cluster_id}/{task_id} [get]
func (h *Handlers) GetPersistentCache(ctx *gin.Context) {
	var params types.PersistentCacheParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	persistentCache, err := h.service.GetPersistentCache(ctx.Request.Context(), params)
	if err != nil {
		ctx.Error(err) // nolint: errcheck
		return
	}

	ctx.JSON(http.StatusOK, persistentCache)
}

// @Summary Get PersistentCaches
// @Description Get PersistentCaches
// @Tags PersistentCache
// @Accept json
// @Produce json
// @Param page query int true "current page" default(0)
// @Param per_page query int true "return max item count, default 10, max 50" default(10) minimum(2) maximum(50)
// @Success 200 {object} []types.GetPersistentCacheResponse
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /api/v1/persistent-caches [get]
func (h *Handlers) GetPersistentCaches(ctx *gin.Context) {
	var query types.GetPersistentCachesQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	h.setPaginationDefault(&query.Page, &query.PerPage)
	persistentCaches, count, err := h.service.GetPersistentCaches(ctx.Request.Context(), query)
	if err != nil {
		ctx.Error(err) // nolint: errcheck
		return
	}

	h.setPaginationLinkHeader(ctx, query.Page, query.PerPage, int(count))
	ctx.JSON(http.StatusOK, persistentCaches)
}
