package job

import (
	"log"
	"time"

	"go-demo/internal/service"

	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	s gocron.Scheduler
}

func StartMetricsGenerator(metricsSvc *service.MetricsService, intervalMin int) (*Scheduler, error) {
	s, err := gocron.NewScheduler(gocron.WithLocation(time.Local))
	if err != nil {
		return nil, err
	}
	if intervalMin <= 0 {
		intervalMin = 30
	}
	_, err = s.NewJob(
		gocron.DurationJob(time.Duration(intervalMin)*time.Minute),
		gocron.NewTask(func() {
			n, err := metricsSvc.GenerateForAll()
			if err != nil {
				log.Printf("[metrics-job] error: %v", err)
				return
			}
			if n > 0 {
				log.Printf("[metrics-job] generated %d metrics", n)
			}
		}),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	)
	if err != nil {
		return nil, err
	}
	s.Start()
	log.Printf("[metrics-job] scheduler started, interval=%dmin", intervalMin)
	return &Scheduler{s: s}, nil
}

func (sch *Scheduler) Stop() {
	if sch != nil && sch.s != nil {
		_ = sch.s.Shutdown()
	}
}
