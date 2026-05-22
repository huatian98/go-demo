package repo

import (
	"go-demo/internal/model"

	"gorm.io/gorm"
)

type UserRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db} }

func (r *UserRepo) FindByOpenid(openid string) (*model.User, error) {
	var u model.User
	if err := r.db.Where("openid = ?", openid).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Create(u *model.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepo) GetByID(id uint64) (*model.User, error) {
	var u model.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) SetDefaultClaim(userID, claimID uint64) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		Update("default_claim_id", claimID).Error
}
