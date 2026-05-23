package service

import (
	"errors"
	"fmt"
	"time"

	"go-demo/internal/model"
	"go-demo/internal/repo"

	"gorm.io/gorm"
)

type ClaimService struct {
	db        *gorm.DB
	claimRepo *repo.ClaimRepo
	jarRepo   *repo.JarRepo
	userRepo  *repo.UserRepo
	timelineRepo *repo.TimelineRepo
}

func NewClaimService(db *gorm.DB, c *repo.ClaimRepo, j *repo.JarRepo, u *repo.UserRepo, t *repo.TimelineRepo) *ClaimService {
	return &ClaimService{db: db, claimRepo: c, jarRepo: j, userRepo: u, timelineRepo: t}
}

type CreateClaimInput struct {
	UserID        uint64
	JarID         uint64
	ApplicantName string
	ContactPhone  string
}

func (s *ClaimService) Create(in CreateClaimInput) (*model.Claim, error) {
	jar, err := s.jarRepo.GetByID(in.JarID)
	if err != nil {
		return nil, errors.New("酒坛不存在")
	}
	if jar.Status != "idle" {
		return nil, errors.New("酒坛已被认领")
	}

	claim := &model.Claim{
		ClaimNo:       generateClaimNo(),
		UserID:        in.UserID,
		JarID:         jar.ID,
		CellarID:      jar.CellarID,
		ApplicantName: in.ApplicantName,
		ContactPhone:  in.ContactPhone,
		Price:         1299.00,
		Status:        "pending",
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 乐观锁更新酒坛
		res := tx.Model(&model.WineJar{}).
			Where("id = ? AND status = ? AND version = ?", jar.ID, "idle", jar.Version).
			Updates(map[string]interface{}{
				"status":           "claimed",
				"current_owner_id": in.UserID,
				"version":          gorm.Expr("version + 1"),
				"claimed_at":       time.Now(),
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("酒坛已被认领,请刷新重试")
		}
		return tx.Create(claim).Error
	})
	if err != nil {
		return nil, err
	}
	return claim, nil
}

func (s *ClaimService) MockPay(claimID uint64) error {
	c, err := s.claimRepo.GetByID(claimID)
	if err != nil {
		return errors.New("认领单不存在")
	}
	if c.Status != "pending" {
		return errors.New("认领单状态不可支付")
	}
	now := time.Now()
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 创建支付记录
		payment := &model.Payment{
			ClaimID:    c.ID,
			OutTradeNo: fmt.Sprintf("MOCK%d%d", time.Now().UnixNano(), c.ID),
			Channel:    "mock",
			Amount:     c.Price,
			Status:     "paid",
			PaidAt:     &now,
		}
		if err := tx.Create(payment).Error; err != nil {
			return err
		}
		// 更新认领状态
		if err := tx.Model(&model.Claim{}).Where("id = ?", c.ID).Updates(map[string]interface{}{
			"status":  "paid",
			"paid_at": now,
		}).Error; err != nil {
			return err
		}
		// 酒坛进入 aging
		if err := tx.Model(&model.WineJar{}).Where("id = ?", c.JarID).
			Update("status", "aging").Error; err != nil {
			return err
		}
		// 设为用户默认展示
		if err := tx.Model(&model.User{}).Where("id = ?", c.UserID).
			Update("default_claim_id", c.ID).Error; err != nil {
			return err
		}
		// 创建初始 timeline
		timelines := []model.JarTimeline{
			{JarID: c.JarID, EventType: "claim", Title: "正式认领", Description: "您已成为这坛酒的守护人", HappenedAt: now},
			{JarID: c.JarID, EventType: "ferment", Title: "入窖陈酿", Description: "酒坛搬入古窖,开启慢呼吸的旅程", HappenedAt: now.Add(time.Hour)},
		}
		for i := range timelines {
			if err := tx.Create(&timelines[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *ClaimService) ListByUser(userID uint64) ([]repo.ClaimVO, error) {
	return s.claimRepo.ListByUserWithJar(userID)
}

func (s *ClaimService) GetByID(id uint64) (*model.Claim, error) {
	return s.claimRepo.GetByID(id)
}

func (s *ClaimService) SetDefault(userID, claimID uint64) error {
	c, err := s.claimRepo.GetByID(claimID)
	if err != nil {
		return errors.New("认领单不存在")
	}
	if c.UserID != userID {
		return errors.New("无权操作")
	}
	return s.userRepo.SetDefaultClaim(userID, claimID)
}

func generateClaimNo() string {
	return fmt.Sprintf("CL%s", time.Now().Format("20060102150405"))
}
