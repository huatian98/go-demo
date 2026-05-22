package middleware

import (
	"strings"

	"go-demo/internal/config"
	jwtpkg "go-demo/internal/pkg/jwt"
	"go-demo/internal/pkg/resp"

	"github.com/gin-gonic/gin"
)

const ctxKeyUserID = "uid"
const ctxKeyOpenid = "openid"

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			resp.Unauthorized(c)
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth {
			resp.Unauthorized(c)
			return
		}
		claims, err := jwtpkg.Parse(token, config.C.JWT.Secret)
		if err != nil {
			resp.Unauthorized(c)
			return
		}
		c.Set(ctxKeyUserID, claims.UserID)
		c.Set(ctxKeyOpenid, claims.Openid)
		c.Next()
	}
}

func GetUserID(c *gin.Context) uint64 {
	v, exists := c.Get(ctxKeyUserID)
	if !exists {
		return 0
	}
	return v.(uint64)
}
