package handlers

import (
	"net/http"

	"d7y.io/dragonfly/v2/internal/job"
	"d7y.io/dragonfly/v2/manager/types"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Summary Create Job
// @Description create by json config
// @Tags Job
// @Accept json
// @Produce json
// @Param Job body types.CreateJobRequest true "Job"
// @Success 200 {object} model.Job
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /jobs [post]
func (h *Handlers) CreateJob(ctx *gin.Context) {
	var json types.CreateJobRequest
	if err := ctx.ShouldBindBodyWith(&json, binding.JSON); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	switch json.Type {
	case job.PreheatJob:
		var json types.CreatePreheatJobRequest
		if err := ctx.ShouldBindBodyWith(&json, binding.JSON); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
			return
		}

		job, err := h.service.CreatePreheatJob(ctx.Request.Context(), json)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, job)
	default:
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": "Unknow type"})
	}
}

// @Summary Destroy Job
// @Description Destroy by id
// @Tags Job
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /jobs/{id} [delete]
func (h *Handlers) DestroyJob(ctx *gin.Context) {
	var params types.JobParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	if err := h.service.DestroyJob(ctx.Request.Context(), params.ID); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Update Job
// @Description Update by json config
// @Tags Job
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param Job body types.UpdateJobRequest true "Job"
// @Success 200 {object} model.Job
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /jobs/{id} [patch]
func (h *Handlers) UpdateJob(ctx *gin.Context) {
	var params types.JobParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.Error(err)
		return
	}

	var json types.UpdateJobRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.Error(err)
		return
	}

	job, err := h.service.UpdateJob(ctx.Request.Context(), params.ID, json)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, job)
}

// @Summary Get Job
// @Description Get Job by id
// @Tags Job
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} model.Job
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /jobs/{id} [get]
func (h *Handlers) GetJob(ctx *gin.Context) {
	var params types.JobParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	job, err := h.service.GetJob(ctx.Request.Context(), params.ID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, job)
}

// @Summary Get Jobs
// @Description Get Jobs
// @Tags Job
// @Accept json
// @Produce json
// @Param page query int true "current page" default(0)
// @Param per_page query int true "return max item count, default 10, max 50" default(10) minimum(2) maximum(50)
// @Success 200 {object} []model.Job
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /jobs [get]
func (h *Handlers) GetJobs(ctx *gin.Context) {
	var query types.GetJobsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	h.setPaginationDefault(&query.Page, &query.PerPage)
	jobs, count, err := h.service.GetJobs(ctx.Request.Context(), query)
	if err != nil {
		ctx.Error(err)
		return
	}

	h.setPaginationLinkHeader(ctx, query.Page, query.PerPage, int(count))
	ctx.JSON(http.StatusOK, jobs)
}
