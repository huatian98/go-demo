package repo

import (
	"go-demo/internal/model"

	"gorm.io/gorm"
)

type MetricsRepo struct{ db *gorm.DB }

func NewMetricsRepo(db *gorm.DB) *MetricsRepo { return &MetricsRepo{db} }

func (r *MetricsRepo) Insert(m *model.JarMetrics) error {
	return r.db.Create(m).Error
}

func (r *MetricsRepo) LatestByJar(jarID uint64) (*model.JarMetrics, error) {
	var m model.JarMetrics
	if err := r.db.Where("jar_id = ?", jarID).
		Order("recorded_at DESC").First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MetricsRepo) HistoryByJar(jarID uint64, days int) ([]model.JarMetrics, error) {
	var list []model.JarMetrics
	if err := r.db.Where("jar_id = ? AND recorded_at > NOW() - (? || ' days')::interval", jarID, days).
		Order("recorded_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// CellarEnv 全局窖藏环境(取最新一条 metrics 的均值,给未认领首页用)
func (r *MetricsRepo) CellarEnv() (map[string]interface{}, error) {
	row := r.db.Raw(`
		SELECT
			COALESCE(AVG(cellar_temperature),18.5) AS temp,
			COALESCE(AVG(cellar_humidity),78.0) AS humidity,
			COALESCE(AVG(ph_level),4.5) AS ph
		FROM jar_metrics
		WHERE recorded_at > NOW() - INTERVAL '2 hours'
	`).Row()
	var temp, humidity, ph float64
	if err := row.Scan(&temp, &humidity, &ph); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"cellar_temperature": temp,
		"cellar_humidity":    humidity,
		"ph_level":           ph,
	}, nil
}
