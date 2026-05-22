package service

import (
	"errors"
	"fmt"
	"time"

	"go-demo/internal/config"
	"go-demo/internal/model"
	jwtpkg "go-demo/internal/pkg/jwt"
	"go-demo/internal/repo"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repo.UserRepo
	cfg      *config.JWTConfig
}

func NewAuthService(u *repo.UserRepo, c *config.JWTConfig) *AuthService {
	return &AuthService{userRepo: u, cfg: c}
}

// WxLogin 简化版微信登录:用 code 当 openid 直接换 JWT(演示阶段)
// 上线后替换为真实调用微信接口:code -> session_key + openid
func (s *AuthService) WxLogin(code string) (string, *model.User, error) {
	if code == "" {
		return "", nil, errors.New("code required")
	}
	openid := "demo_" + code
	u, err := s.userRepo.FindByOpenid(openid)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, err
		}
		u = &model.User{
			Openid:   openid,
			Nickname: fmt.Sprintf("酒友_%s", code[:min(len(code), 6)]),
			Avatar:   "/static/avatar-default.png",
		}
		if err := s.userRepo.Create(u); err != nil {
			return "", nil, err
		}
	}
	token, err := jwtpkg.Sign(u.ID, u.Openid, s.cfg.Secret, s.cfg.ExpireHours)
	if err != nil {
		return "", nil, err
	}
	return token, u, nil
}

func (s *AuthService) GetUser(id uint64) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var _ = time.Now
