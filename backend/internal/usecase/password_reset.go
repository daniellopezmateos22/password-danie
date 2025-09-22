// Caso de uso de password reset: generar token de recuperación y confirmar cambio de contraseña.
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

func NewPasswordReset(users repository.UserRepo) *PasswordReset {
	return &PasswordReset{users: users}
}

// Request genera un token temporal y lo guarda en el usuario.
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

	// actualizar en BD (requiere UpdateReset en UserRepo)
	if err := pr.users.UpdateReset(u.ID, &token, &exp); err != nil {
		return "", err
	}
	return token, nil
}

// Confirm valida el token y actualiza el password.
func (pr *PasswordReset) Confirm(token, newPassword string) error {
	u, err := pr.users.GetByResetToken(token)
	if err != nil {
		return err
	}
	if u == nil || u.ResetExpiresAt == nil || time.Now().After(*u.ResetExpiresAt) {
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
