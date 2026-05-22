package handler

import (
	"go-demo/internal/middleware"
	"go-demo/internal/pkg/resp"
	"go-demo/internal/service"

	"github.com/gin-gonic/gin"
)

type ClaimHandler struct {
	svc *service.ClaimService
}

func NewClaimHandler(s *service.ClaimService) *ClaimHandler { return &ClaimHandler{s} }

type createClaimReq struct {
	JarID         uint64 `json:"jar_id" binding:"required"`
	ApplicantName string `json:"applicant_name" binding:"required"`
	ContactPhone  string `json:"contact_phone" binding:"required"`
}

// POST /api/v1/claims
func (h *ClaimHandler) Create(c *gin.Context) {
	uid := middleware.GetUserID(c)
	var req createClaimReq
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}
	claim, err := h.svc.Create(service.CreateClaimInput{
		UserID:        uid,
		JarID:         req.JarID,
		ApplicantName: req.ApplicantName,
		ContactPhone:  req.ContactPhone,
	})
	if err != nil {
		resp.Fail(c, 400, err.Error())
		return
	}
	resp.OK(c, claim)
}

// GET /api/v1/claims
func (h *ClaimHandler) List(c *gin.Context) {
	uid := middleware.GetUserID(c)
	list, err := h.svc.ListByUser(uid)
	if err != nil {
		resp.ServerError(c, err.Error())
		return
	}
	resp.OK(c, list)
}

// GET /api/v1/claims/:id
func (h *ClaimHandler) Detail(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		resp.BadRequest(c, "invalid id")
		return
	}
	claim, err := h.svc.GetByID(id)
	if err != nil {
		resp.NotFound(c, "claim not found")
		return
	}
	resp.OK(c, claim)
}

// POST /api/v1/claims/:id/set-default
func (h *ClaimHandler) SetDefault(c *gin.Context) {
	uid := middleware.GetUserID(c)
	id, err := parseID(c.Param("id"))
	if err != nil {
		resp.BadRequest(c, "invalid id")
		return
	}
	if err := h.svc.SetDefault(uid, id); err != nil {
		resp.Fail(c, 400, err.Error())
		return
	}
	resp.OK(c, gin.H{"ok": true})
}

type mockPayReq struct {
	ClaimID uint64 `json:"claim_id" binding:"required"`
}

// POST /api/v1/payments/mock-pay
func (h *ClaimHandler) MockPay(c *gin.Context) {
	var req mockPayReq
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.MockPay(req.ClaimID); err != nil {
		resp.Fail(c, 400, err.Error())
		return
	}
	resp.OK(c, gin.H{"paid": true, "claim_id": req.ClaimID})
}
