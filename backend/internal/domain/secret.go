// Package domain define entidades del dominio. Secret modela una credencial cifrada del vault.
package domain

import "time"

type Secret struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	Username       string    `json:"username"`
	URL            string    `json:"url"`
	URLDomain      string    `json:"url_domain"` 
	Notes          string    `json:"notes"`
	Icon           string    `json:"icon"`
	Title          string    `json:"title"`      
	PasswordCipher string    `json:"-"`        
	PasswordIV     string    `json:"-"`        
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
