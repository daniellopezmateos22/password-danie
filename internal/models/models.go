package models

import "time"

// Usuarios
type User struct {
	ID           uint      `gorm:"primaryKey"`
	Email        string    `gorm:"uniqueIndex"`
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Ítems del vault 
type VaultItem struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"index"`
	Title      string
	Username   string
	PasswordC  []byte
	URL        string
	NotesC     []byte
	Icon       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Tokens de reset
type ResetToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"`
	Token     string    `gorm:"uniqueIndex"`
	ExpiresAt time.Time
	Used      bool      `gorm:"default:false"`
	CreatedAt time.Time
}
