// Package repository declara puertos (interfaces) de persistencia para usuarios.
package repository

import (
	"time"

	"password-danie/internal/domain"
)

type UserRepo interface {
	// básicos
	Create(email, passwordHash string) (int64, error)
	GetByID(id int64) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)

	// password reset
	UpdateReset(userID int64, token *string, expiresAt *time.Time) error
	GetByResetToken(token string) (*domain.User, error)
	UpdatePassword(userID int64, passwordHash string) error
}
