package handler

import (
	"go-demo/internal/middleware"
	"go-demo/internal/pkg/resp"
	"go-demo/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler { return &AuthHandler{s} }

type wxLoginReq struct {
	Code string `json:"code"`
}

// POST /api/v1/auth/wx-login
func (h *AuthHandler) WxLogin(c *gin.Context) {
	var req wxLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "invalid request")
		return
	}
	if req.Code == "" {
		req.Code = "default"
	}
	token, user, err := h.svc.WxLogin(req.Code)
	if err != nil {
		resp.ServerError(c, err.Error())
		return
	}
	resp.OK(c, gin.H{"token": token, "user": user})
}

// GET /api/v1/user/me
func (h *AuthHandler) Me(c *gin.Context) {
	uid := middleware.GetUserID(c)
	u, err := h.svc.GetUser(uid)
	if err != nil {
		resp.NotFound(c, "user not found")
		return
	}
	resp.OK(c, u)
}
