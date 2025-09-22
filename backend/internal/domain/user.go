// Package domain define entidades del dominio. User incluye soporte de reset de contraseña.
package domain

import "time"

type User struct {
	ID             int64     `json:"id"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"` 
	ResetToken     string    `json:"-"` 
	ResetExpiresAt time.Time `json:"-"` 
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
