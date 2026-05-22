package router

import (
	"go-demo/internal/handler"
	"go-demo/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth   *handler.AuthHandler
	Home   *handler.HomeHandler
	Jar    *handler.JarHandler
	Claim  *handler.ClaimHandler
	Legacy *handler.LegacyHandler
}

func Setup(h *Handlers) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	r.GET("/", func(c *gin.Context) {
		c.String(200, "时光酿 API · 服务已启动")
	})
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	// 兼容老接口
	r.GET("/api/claim/:id", h.Legacy.ClaimByJarID)

	v1 := r.Group("/api/v1")
	{
		// 公开接口
		v1.POST("/auth/wx-login", h.Auth.WxLogin)
		v1.GET("/series", h.Home.SeriesList)
		v1.GET("/components", h.Home.Components)
		v1.GET("/home/cellar-env", h.Home.CellarEnv)
		v1.GET("/home/craft-steps", h.Home.CraftSteps)
		v1.GET("/jars/:id/metrics/latest", h.Jar.MetricsLatest)
		v1.GET("/jars/:id/metrics/history", h.Jar.MetricsHistory)
		v1.GET("/jars/:id/timeline", h.Home.Timeline)

		// 鉴权接口
		auth := v1.Group("")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/user/me", h.Auth.Me)
			auth.GET("/home/dashboard", h.Home.Dashboard)
			auth.POST("/claims", h.Claim.Create)
			auth.GET("/claims", h.Claim.List)
			auth.GET("/claims/:id", h.Claim.Detail)
			auth.POST("/claims/:id/set-default", h.Claim.SetDefault)
			auth.POST("/payments/mock-pay", h.Claim.MockPay)
		}
	}

	return r
}
