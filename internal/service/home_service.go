package service

import (
	"errors"

	"go-demo/internal/model"
	"go-demo/internal/repo"

	"gorm.io/gorm"
)

type HomeService struct {
	userRepo     *repo.UserRepo
	claimRepo    *repo.ClaimRepo
	jarRepo      *repo.JarRepo
	cellarRepo   *repo.CellarRepo
	seriesRepo   *repo.SeriesRepo
	metricsRepo  *repo.MetricsRepo
	timelineRepo *repo.TimelineRepo
	contentRepo  *repo.ContentRepo
}

func NewHomeService(
	u *repo.UserRepo, c *repo.ClaimRepo, j *repo.JarRepo, cl *repo.CellarRepo,
	s *repo.SeriesRepo, m *repo.MetricsRepo, t *repo.TimelineRepo, ct *repo.ContentRepo,
) *HomeService {
	return &HomeService{
		userRepo: u, claimRepo: c, jarRepo: j, cellarRepo: cl,
		seriesRepo: s, metricsRepo: m, timelineRepo: t, contentRepo: ct,
	}
}

type CellarEnvResp struct {
	InCellarTemp      float64           `json:"in_cellar_temp"`
	InCellarHumidity  float64           `json:"in_cellar_humidity"`
	OutCellarTemp     float64           `json:"out_cellar_temp"`
	OutCellarHumidity float64           `json:"out_cellar_humidity"`
	WinePh            float64           `json:"wine_ph"`
	CraftSteps        []model.CraftStep `json:"craft_steps"`
}

func (s *HomeService) CellarEnv() (*CellarEnvResp, error) {
	env, err := s.metricsRepo.CellarEnv()
	if err != nil {
		return nil, err
	}
	steps, _ := s.contentRepo.ListCraftSteps()
	return &CellarEnvResp{
		InCellarTemp:      env["in_cellar_temp"].(float64),
		InCellarHumidity:  env["in_cellar_humidity"].(float64),
		OutCellarTemp:     env["out_cellar_temp"].(float64),
		OutCellarHumidity: env["out_cellar_humidity"].(float64),
		WinePh:            env["wine_ph"].(float64),
		CraftSteps:        steps,
	}, nil
}

type DashboardResp struct {
	State      string                  `json:"state"` // not_claimed / claimed
	Claim      *model.Claim            `json:"claim,omitempty"`
	Jar        *model.WineJar          `json:"jar,omitempty"`
	Series     *model.WineSeries       `json:"series,omitempty"`
	Cellar     *model.Cellar           `json:"cellar,omitempty"`
	Metrics    *model.JarMetrics       `json:"metrics,omitempty"`
	Timelines  []model.JarTimeline     `json:"timelines,omitempty"`
	Components []model.WineComponent   `json:"components,omitempty"`
	AgingDays  int                     `json:"aging_days"`
}

func (s *HomeService) Dashboard(userID uint64, metricsSvc *MetricsService) (*DashboardResp, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user.DefaultClaimID == 0 {
		return &DashboardResp{State: "not_claimed"}, nil
	}
	claim, err := s.claimRepo.GetByID(user.DefaultClaimID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &DashboardResp{State: "not_claimed"}, nil
		}
		return nil, err
	}
	jar, _ := s.jarRepo.GetByID(claim.JarID)
	series, _ := s.seriesRepo.GetByID(jar.SeriesID)
	cellar, _ := s.cellarRepo.GetByID(jar.CellarID)
	metrics, _ := metricsSvc.Latest(jar.ID)
	timelines, _ := s.timelineRepo.ListByJar(jar.ID)
	components, _ := s.contentRepo.ListComponents()

	agingDays := 0
	if jar != nil && jar.ClaimedAt != nil {
		agingDays = int(metrics.RecordedAt.Sub(*jar.ClaimedAt).Hours() / 24)
	}

	return &DashboardResp{
		State:      "claimed",
		Claim:      claim,
		Jar:        jar,
		Series:     series,
		Cellar:     cellar,
		Metrics:    metrics,
		Timelines:  timelines,
		Components: components,
		AgingDays:  agingDays,
	}, nil
}
