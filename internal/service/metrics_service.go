package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"go-demo/internal/model"
	"go-demo/internal/repo"
)

type MetricsService struct {
	jarRepo     *repo.JarRepo
	metricsRepo *repo.MetricsRepo
}

func NewMetricsService(j *repo.JarRepo, m *repo.MetricsRepo) *MetricsService {
	return &MetricsService{jarRepo: j, metricsRepo: m}
}

// Latest 通过 jar 内部 ID 查找(向后兼容老 handler)
func (s *MetricsService) Latest(jarInternalID uint64) (*model.JarMetrics, error) {
	jar, err := s.jarRepo.GetByID(jarInternalID)
	if err != nil {
		return nil, err
	}
	m, err := s.metricsRepo.LatestByJar(jar.Code)
	if err != nil {
		// 没有数据时即时合成一条(不入库)
		return s.synthesize(jar, false), nil
	}
	return m, nil
}

// LatestByCode 直接用酒坛编号查
func (s *MetricsService) LatestByCode(wineJarID string) (*model.JarMetrics, error) {
	m, err := s.metricsRepo.LatestByJar(wineJarID)
	if err != nil {
		jar, ferr := s.jarRepo.GetByCode(wineJarID)
		if ferr != nil {
			return nil, ferr
		}
		return s.synthesize(jar, false), nil
	}
	return m, nil
}

func (s *MetricsService) History(jarInternalID uint64, days int) ([]model.JarMetrics, error) {
	if days <= 0 || days > 30 {
		days = 7
	}
	jar, err := s.jarRepo.GetByID(jarInternalID)
	if err != nil {
		return nil, err
	}
	return s.metricsRepo.HistoryByJar(jar.Code, days)
}

func (s *MetricsService) CellarEnv() (map[string]interface{}, error) {
	return s.metricsRepo.CellarEnv()
}

// GenerateForAll 给所有 active 酒坛生成一条 metrics(供定时任务调用)
func (s *MetricsService) GenerateForAll() (int, error) {
	jars, err := s.jarRepo.ListActive()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, jar := range jars {
		m := s.synthesize(&jar, true)
		if err := s.metricsRepo.Insert(m); err != nil {
			continue
		}
		count++
	}
	return count, nil
}

// synthesize 根据酒坛状态合成一条 metrics
func (s *MetricsService) synthesize(jar *model.WineJar, persist bool) *model.JarMetrics {
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(jar.ID)))

	agingDays := 0
	if jar.ClaimedAt != nil {
		agingDays = int(time.Since(*jar.ClaimedAt).Hours() / 24)
	}

	inTemp := 18.0 + rng.Float64()*2.0    // 18~20℃
	inHum := 75.0 + rng.Float64()*5.0     // 75~80%
	outTemp := 20.0 + rng.Float64()*8.0   // 20~28℃
	outHum := 55.0 + rng.Float64()*20.0   // 55~75%
	ph := 4.3 + rng.Float64()*0.5         // 4.3~4.8

	phStatus := "稳定"
	if ph > 4.7 {
		phStatus = "偏高"
	}

	breathing := pickBreathing(agingDays)
	narrative := generateNarrative(jar.Code, ph, inTemp, breathing)

	return &model.JarMetrics{
		WineJarID:         jar.Code,
		WinePh:            round2(ph),
		PhStatus:          phStatus,
		InCellarTemp:      round1(inTemp),
		InCellarHumidity:  round1(inHum),
		OutCellarTemp:     round1(outTemp),
		OutCellarHumidity: round1(outHum),
		BreathingState:    breathing,
		AINarrative:       narrative,
		RecordedAt:        time.Now(),
	}
}

func pickBreathing(agingDays int) string {
	switch {
	case agingDays < 30:
		return "初醒发酵中"
	case agingDays < 90:
		return "活跃发酵中"
	case agingDays < 180:
		return "风味沉淀中"
	case agingDays < 365:
		return "深度陈酿中"
	default:
		return "沉睡养护中"
	}
}

func generateNarrative(code string, ph, temp float64, state string) string {
	if ph > 4.7 {
		return fmt.Sprintf("酸度略高(%.2f),%s 的脸颊微微泛红,需注意窖温调控", ph, code)
	}
	if temp > 19.5 {
		return fmt.Sprintf("窖温偏暖(%.1f°C),%s 略显躁动,菌群活跃度上升", temp, code)
	}
	return fmt.Sprintf("当前\"红曲之灵\"正处于舒适的%s状态,呼吸均匀,风味平稳积累中", state)
}

func round1(v float64) float64 { return float64(int(v*10)) / 10 }
func round2(v float64) float64 { return float64(int(v*100)) / 100 }

var ErrNoJar = errors.New("jar not found")
