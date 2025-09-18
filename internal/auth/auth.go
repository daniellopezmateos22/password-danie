package auth

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(getEnv("JWT_SECRET", "dev-secret"))

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func MakeToken(userID uint) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	return t.SignedString(jwtSecret)
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "token requerido"})
			return
		}
		raw := strings.TrimPrefix(h, "Bearer ")
		parsed, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("alg inválido")
			}
			return jwtSecret, nil
		})
		if err != nil || !parsed.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "token inválido"})
			return
		}
		claims := parsed.Claims.(jwt.MapClaims)
		sub, ok := claims["sub"].(float64)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "claims inválidos"})
			return
		}
		c.Set("user_id", uint(sub))
		c.Next()
	}
}
