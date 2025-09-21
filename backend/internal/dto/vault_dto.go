// Package dto contiene structs de petición/respuesta para el módulo de vault.
package dto

type CreateSecretRequest struct {
	Username      string  `json:"username" binding:"required"`
	PasswordPlain string  `json:"password_plain" binding:"required"`
	URL           string  `json:"url"`
	Notes         string  `json:"notes"`
	Icon          string  `json:"icon"`
	Title         *string `json:"title"` 
}

type UpdateSecretRequest struct {
	Username      *string `json:"username"`
	PasswordPlain *string `json:"password_plain"`
	URL           *string `json:"url"`
	Notes         *string `json:"notes"`
	Icon          *string `json:"icon"`
	Title         *string `json:"title"`
}
