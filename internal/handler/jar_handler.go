package handler

import (
	"strconv"

	"go-demo/internal/pkg/resp"
	"go-demo/internal/repo"
	"go-demo/internal/service"

	"github.com/gin-gonic/gin"
)

type JarHandler struct {
	metricsSvc *service.MetricsService
	jarRepo    *repo.JarRepo
}

type JarHandlerDeps struct {
	MetricsSvc *service.MetricsService
	JarRepo    *repo.JarRepo
}

func NewJarHandler(d JarHandlerDeps) *JarHandler {
	return &JarHandler{metricsSvc: d.MetricsSvc, jarRepo: d.JarRepo}
}

// GET /api/v1/jars/available
func (h *JarHandler) Available(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	list, err := h.jarRepo.ListAvailableBrief(limit)
	if err != nil {
		resp.ServerError(c, err.Error())
		return
	}
	if list == nil {
		list = []repo.JarBrief{}
	}
	resp.OK(c, list)
}

// GET /api/v1/jars/:id/metrics/latest
func (h *JarHandler) MetricsLatest(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		resp.BadRequest(c, "invalid id")
		return
	}
	m, err := h.metricsSvc.Latest(id)
	if err != nil {
		resp.NotFound(c, "no metrics")
		return
	}
	resp.OK(c, m)
}

// GET /api/v1/jars/:id/metrics/history?days=7
func (h *JarHandler) MetricsHistory(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		resp.BadRequest(c, "invalid id")
		return
	}
	days, _ := strconv.Atoi(c.Query("days"))
	list, err := h.metricsSvc.History(id, days)
	if err != nil {
		resp.ServerError(c, err.Error())
		return
	}
	resp.OK(c, list)
}

func parseID(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}
