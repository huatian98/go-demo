package repo

import (
	"go-demo/internal/model"

	"gorm.io/gorm"
)

type ClaimRepo struct{ db *gorm.DB }

func NewClaimRepo(db *gorm.DB) *ClaimRepo { return &ClaimRepo{db} }

func (r *ClaimRepo) Create(claim *model.Claim) error {
	return r.db.Create(claim).Error
}

func (r *ClaimRepo) GetByID(id uint64) (*model.Claim, error) {
	var c model.Claim
	if err := r.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ClaimRepo) ListByUser(userID uint64) ([]model.Claim, error) {
	var list []model.Claim
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ClaimRepo) UpdateStatus(id uint64, status string) error {
	return r.db.Model(&model.Claim{}).Where("id = ?", id).Update("status", status).Error
}

type TimelineRepo struct{ db *gorm.DB }

func NewTimelineRepo(db *gorm.DB) *TimelineRepo { return &TimelineRepo{db} }

func (r *TimelineRepo) ListByJar(jarID uint64) ([]model.JarTimeline, error) {
	var list []model.JarTimeline
	if err := r.db.Where("jar_id = ?", jarID).
		Order("happened_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *TimelineRepo) Insert(t *model.JarTimeline) error {
	return r.db.Create(t).Error
}

type SeriesRepo struct{ db *gorm.DB }

func NewSeriesRepo(db *gorm.DB) *SeriesRepo { return &SeriesRepo{db} }

func (r *SeriesRepo) List() ([]model.WineSeries, error) {
	var list []model.WineSeries
	if err := r.db.Where("status = ?", 1).Order("sort ASC, id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *SeriesRepo) GetByID(id uint64) (*model.WineSeries, error) {
	var s model.WineSeries
	if err := r.db.First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

type CellarRepo struct{ db *gorm.DB }

func NewCellarRepo(db *gorm.DB) *CellarRepo { return &CellarRepo{db} }

func (r *CellarRepo) GetByID(id uint64) (*model.Cellar, error) {
	var c model.Cellar
	if err := r.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CellarRepo) List() ([]model.Cellar, error) {
	var list []model.Cellar
	if err := r.db.Order("id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

type ContentRepo struct{ db *gorm.DB }

func NewContentRepo(db *gorm.DB) *ContentRepo { return &ContentRepo{db} }

func (r *ContentRepo) ListComponents() ([]model.WineComponent, error) {
	var list []model.WineComponent
	if err := r.db.Order("sort ASC, id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ContentRepo) ListCraftSteps() ([]model.CraftStep, error) {
	var list []model.CraftStep
	if err := r.db.Order("step_no ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
