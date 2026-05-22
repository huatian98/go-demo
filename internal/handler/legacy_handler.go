package handler

import (
	"go-demo/internal/pkg/resp"
	"go-demo/internal/repo"

	"github.com/gin-gonic/gin"
)

// LegacyHandler 兼容老版本接口(小程序 claim.js 用的 /api/claim/:id 旧 schema)
type LegacyHandler struct {
	jarRepo    *repo.JarRepo
	seriesRepo *repo.SeriesRepo
	cellarRepo *repo.CellarRepo
}

func NewLegacyHandler(j *repo.JarRepo, s *repo.SeriesRepo, c *repo.CellarRepo) *LegacyHandler {
	return &LegacyHandler{j, s, c}
}

// GET /api/claim/:id
// 旧接口:返回 { code, series, applicant, phone, cellar, address }
// 实际从数据库读取 jar + series + cellar 拼装
func (h *LegacyHandler) ClaimByJarID(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		resp.BadRequest(c, "invalid id")
		return
	}
	jar, err := h.jarRepo.GetByID(id)
	if err != nil {
		resp.NotFound(c, "not found")
		return
	}
	series, _ := h.seriesRepo.GetByID(jar.SeriesID)
	cellar, _ := h.cellarRepo.GetByID(jar.CellarID)

	seriesName := ""
	if series != nil {
		seriesName = series.Name
	}
	cellarName, address := "", ""
	if cellar != nil {
		cellarName = cellar.Name
		address = cellar.Address
	}

	resp.OK(c, gin.H{
		"code":      jar.Code,
		"series":    seriesName,
		"applicant": "可乐",
		"phone":     "138 **** 5678",
		"cellar":    cellarName,
		"address":   address,
	})
}
