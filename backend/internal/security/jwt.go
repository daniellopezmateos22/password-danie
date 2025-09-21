// Utilidades JWT (HS256): generar tokens de acceso y validar/extraer claims.
package security

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func jwtKey() []byte {
	k := os.Getenv("JWT_SECRET")
	if k == "" {
		k = "dev-secret-change-me"
	}
	return []byte(k)
}

func GenerateAccessToken(sub int64, email string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   sub,
		"email": email,
		"exp":   time.Now().Add(ttl).Unix(),
	})
	return token.SignedString(jwtKey())
}

func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey(), nil
	})
	if err != nil || !t.Valid {
		if err == nil {
			err = errors.New("invalid token")
		}
		return nil, err
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}
