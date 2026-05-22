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

// LatestByJar 通过酒坛编号(wine_jar_id,即 wine_jars.code)查最新一条
func (r *MetricsRepo) LatestByJar(wineJarID string) (*model.JarMetrics, error) {
	var m model.JarMetrics
	if err := r.db.Where("wine_jar_id = ?", wineJarID).
		Order("recorded_at DESC").First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MetricsRepo) HistoryByJar(wineJarID string, days int) ([]model.JarMetrics, error) {
	var list []model.JarMetrics
	if err := r.db.Where("wine_jar_id = ? AND recorded_at > NOW() - (? || ' days')::interval", wineJarID, days).
		Order("recorded_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// CellarEnv 全局窖藏环境(取最新一段时间内 metrics 的均值,给未认领首页用)
func (r *MetricsRepo) CellarEnv() (map[string]interface{}, error) {
	row := r.db.Raw(`
		SELECT
			COALESCE(AVG(in_cellar_temp), 18.5)     AS in_temp,
			COALESCE(AVG(in_cellar_humidity), 78.0) AS in_humidity,
			COALESCE(AVG(out_cellar_temp), 24.0)    AS out_temp,
			COALESCE(AVG(out_cellar_humidity), 65.0) AS out_humidity,
			COALESCE(AVG(wine_ph), 4.5)             AS wine_ph
		FROM jar_metrics
		WHERE recorded_at > NOW() - INTERVAL '2 hours'
	`).Row()
	var inTemp, inHum, outTemp, outHum, ph float64
	if err := row.Scan(&inTemp, &inHum, &outTemp, &outHum, &ph); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"in_cellar_temp":      inTemp,
		"in_cellar_humidity":  inHum,
		"out_cellar_temp":     outTemp,
		"out_cellar_humidity": outHum,
		"wine_ph":             ph,
	}, nil
}
