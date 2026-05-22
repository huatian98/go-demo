package repo

import (
	"go-demo/internal/model"

	"gorm.io/gorm"
)

type JarRepo struct{ db *gorm.DB }

func NewJarRepo(db *gorm.DB) *JarRepo { return &JarRepo{db} }

func (r *JarRepo) GetByID(id uint64) (*model.WineJar, error) {
	var j model.WineJar
	if err := r.db.First(&j, id).Error; err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *JarRepo) GetByCode(code string) (*model.WineJar, error) {
	var j model.WineJar
	if err := r.db.Where("code = ?", code).First(&j).Error; err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *JarRepo) ListAvailable(limit int) ([]model.WineJar, error) {
	var jars []model.WineJar
	if err := r.db.Where("status = ?", "idle").Limit(limit).Find(&jars).Error; err != nil {
		return nil, err
	}
	return jars, nil
}

// JarBrief 酒坛 + 系列 + 酒窖关联,给选坛页用
type JarBrief struct {
	ID         uint64  `json:"id"`
	Code       string  `json:"code"`
	Year       int     `json:"year"`
	CoverURL   string  `json:"cover_url"`
	Status     string  `json:"status"`
	SeriesID   uint64  `json:"series_id"`
	SeriesName string  `json:"series_name"`
	BasePrice  float64 `json:"base_price"`
	CellarID   uint64  `json:"cellar_id"`
	CellarName string  `json:"cellar_name"`
	Address    string  `json:"address"`
}

func (r *JarRepo) ListAvailableBrief(limit int) ([]JarBrief, error) {
	if limit <= 0 {
		limit = 50
	}
	var briefs []JarBrief
	err := r.db.Table("wine_jars").
		Select(`wine_jars.id, wine_jars.code, wine_jars.year, wine_jars.cover_url, wine_jars.status,
		        wine_series.id  AS series_id, wine_series.name AS series_name, wine_series.base_price,
		        cellars.id      AS cellar_id, cellars.name AS cellar_name, cellars.address`).
		Joins("LEFT JOIN wine_series ON wine_series.id = wine_jars.series_id").
		Joins("LEFT JOIN cellars ON cellars.id = wine_jars.cellar_id").
		Where("wine_jars.status = ?", "idle").
		Order("wine_jars.id ASC").
		Limit(limit).
		Scan(&briefs).Error
	return briefs, err
}

// ListActive 返回所有已认领且在养护中的酒坛(给定时任务用)
func (r *JarRepo) ListActive() ([]model.WineJar, error) {
	var jars []model.WineJar
	if err := r.db.Where("status IN ?", []string{"claimed", "aging", "ready"}).Find(&jars).Error; err != nil {
		return nil, err
	}
	return jars, nil
}

// ClaimWithLock 用乐观锁认领酒坛
func (r *JarRepo) ClaimWithLock(jarID, ownerID uint64, expectedVersion int) (int64, error) {
	tx := r.db.Model(&model.WineJar{}).
		Where("id = ? AND status = ? AND version = ?", jarID, "idle", expectedVersion).
		Updates(map[string]interface{}{
			"status":           "claimed",
			"current_owner_id": ownerID,
			"version":          gormExpr("version + 1"),
		})
	return tx.RowsAffected, tx.Error
}

func gormExpr(s string) interface{} {
	return gorm.Expr(s)
}
