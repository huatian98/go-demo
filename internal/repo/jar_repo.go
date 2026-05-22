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
