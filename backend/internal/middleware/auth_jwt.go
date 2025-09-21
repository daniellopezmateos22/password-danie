// Middleware Gin de autenticación JWT: exige Bearer token y expone claims y userID en contexto.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"password-danie/internal/security"
)

const CtxClaims = "claims"
const CtxUserID = "userID"

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		token := strings.TrimPrefix(h, "Bearer ")
		claims, err := security.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set(CtxClaims, claims)

		var uid int64
		if v, ok := claims["sub"]; ok {
			switch t := v.(type) {
			case float64:
				uid = int64(t)
			case int64:
				uid = t
			case int:
				uid = int64(t)
			}
		}
		if uid == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid subject"})
			return
		}
		c.Set(CtxUserID, uid)
		c.Next()
	}
}
