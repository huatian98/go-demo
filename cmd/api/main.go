package main

import (
	"fmt"
	"log"

	"go-demo/internal/config"
	"go-demo/internal/handler"
	"go-demo/internal/job"
	"go-demo/internal/repo"
	"go-demo/internal/router"
	"go-demo/internal/service"
)

func main() {
	// 加载配置
	c := config.Load()

	// 初始化 DB
	db := repo.InitDB(&c.Database)

	// repo 层
	userRepo := repo.NewUserRepo(db)
	jarRepo := repo.NewJarRepo(db)
	metricsRepo := repo.NewMetricsRepo(db)
	claimRepo := repo.NewClaimRepo(db)
	timelineRepo := repo.NewTimelineRepo(db)
	seriesRepo := repo.NewSeriesRepo(db)
	cellarRepo := repo.NewCellarRepo(db)
	contentRepo := repo.NewContentRepo(db)

	// service 层
	authSvc := service.NewAuthService(userRepo, &c.JWT)
	metricsSvc := service.NewMetricsService(jarRepo, metricsRepo)
	claimSvc := service.NewClaimService(db, claimRepo, jarRepo, userRepo, timelineRepo)
	homeSvc := service.NewHomeService(userRepo, claimRepo, jarRepo, cellarRepo, seriesRepo, metricsRepo, timelineRepo, contentRepo)

	// handler 层
	handlers := &router.Handlers{
		Auth:   handler.NewAuthHandler(authSvc),
		Home:   handler.NewHomeHandler(homeSvc, metricsSvc, seriesRepo, contentRepo, timelineRepo),
		Jar:    handler.NewJarHandler(handler.JarHandlerDeps{MetricsSvc: metricsSvc, JarRepo: jarRepo}),
		Claim:  handler.NewClaimHandler(claimSvc),
		Legacy: handler.NewLegacyHandler(jarRepo, seriesRepo, cellarRepo),
	}

	// 启动定时任务
	sch, err := job.StartMetricsGenerator(metricsSvc, c.Metrics.IntervalMinutes)
	if err != nil {
		log.Fatalf("start metrics job failed: %v", err)
	}
	defer sch.Stop()

	// 启动 HTTP 服务
	r := router.Setup(handlers)
	addr := fmt.Sprintf(":%d", c.App.Port)
	log.Printf("server listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("run server failed: %v", err)
	}
}
