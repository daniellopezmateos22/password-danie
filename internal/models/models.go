package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type VaultItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"-"`
	Title     string    `gorm:"not null" json:"title"`
	Username  string    `gorm:"not null" json:"username"`
	PasswordC string    `gorm:"not null" json:"-"` 
	URL       string    `json:"url"`
	NotesC    string    `json:"-"`                 
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResetToken struct {
	ID        uint      `gorm:"primaryKey" json:"-"`
	UserID    uint      `gorm:"index;not null" json:"-"`
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"-"`
}
