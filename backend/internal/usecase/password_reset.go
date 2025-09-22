// Caso de uso de password reset: generar token y confirmar con valores (no punteros) en el dominio.
package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"password-danie/internal/repository"
)

type PasswordReset struct {
	users repository.UserRepo
}

func NewPasswordReset(users repository.UserRepo) *PasswordReset { return &PasswordReset{users: users} }

func (pr *PasswordReset) Request(email string) (string, error) {
	u, err := pr.users.GetByEmail(email)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", errors.New("user not found")
	}
	token := randomToken(32)
	exp := time.Now().Add(1 * time.Hour)
	if err := pr.users.UpdateReset(u.ID, &token, &exp); err != nil {
		return "", err
	}
	return token, nil
}

func (pr *PasswordReset) Confirm(token, newPassword string) error {
	u, err := pr.users.GetByResetToken(token)
	if err != nil {
		return err
	}
	// Validar existencia y expiración usando valores
	if u == nil || u.ResetToken == "" || u.ResetExpiresAt.IsZero() || time.Now().After(u.ResetExpiresAt) {
		return errors.New("invalid or expired token")
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	return pr.users.UpdatePassword(u.ID, string(hash))
}

func randomToken(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
