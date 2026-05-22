package handler

import (
	"go-demo/internal/middleware"
	"go-demo/internal/pkg/resp"
	"go-demo/internal/repo"
	"go-demo/internal/service"

	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	homeSvc      *service.HomeService
	metricsSvc   *service.MetricsService
	seriesRepo   *repo.SeriesRepo
	contentRepo  *repo.ContentRepo
	timelineRepo *repo.TimelineRepo
}

func NewHomeHandler(
	hs *service.HomeService,
	ms *service.MetricsService,
	sr *repo.SeriesRepo,
	cr *repo.ContentRepo,
	tr *repo.TimelineRepo,
) *HomeHandler {
	return &HomeHandler{homeSvc: hs, metricsSvc: ms, seriesRepo: sr, contentRepo: cr, timelineRepo: tr}
}

// GET /api/v1/home/cellar-env
func (h *HomeHandler) CellarEnv(c *gin.Context) {
	data, err := h.homeSvc.CellarEnv()
	if err != nil {
		resp.ServerError(c, err.Error())
		return
	}
	resp.OK(c, data)
}

// GET /api/v1/home/dashboard
func (h *HomeHandler) Dashboard(c *gin.Context) {
	uid := middleware.GetUserID(c)
	data, err := h.homeSvc.Dashboard(uid, h.metricsSvc)
	if err != nil {
		resp.ServerError(c, err.Error())
		return
	}
	resp.OK(c, data)
}

// GET /api/v1/home/craft-steps
func (h *HomeHandler) CraftSteps(c *gin.Context) {
	steps, err := h.contentRepo.ListCraftSteps()
	if err != nil {
		resp.ServerError(c, err.Error())
		return
	}
	resp.OK(c, steps)
}

// GET /api/v1/components
func (h *HomeHandler) Components(c *gin.Context) {
	list, err := h.contentRepo.ListComponents()
	if err != nil {
		resp.ServerError(c, err.Error())
		return
	}
	resp.OK(c, list)
}

// GET /api/v1/series
func (h *HomeHandler) SeriesList(c *gin.Context) {
	list, err := h.seriesRepo.List()
	if err != nil {
		resp.ServerError(c, err.Error())
		return
	}
	resp.OK(c, list)
}

// GET /api/v1/jars/:id/timeline
func (h *HomeHandler) Timeline(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		resp.BadRequest(c, "invalid id")
		return
	}
	list, err := h.timelineRepo.ListByJar(id)
	if err != nil {
		resp.ServerError(c, err.Error())
		return
	}
	resp.OK(c, list)
}
