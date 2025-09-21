// Caso de uso del vault: CRUD de secretos con cifrado AES-GCM y soporte de búsqueda/filtrado.
package usecase

import (
	"errors"

	"password-danie/internal/domain"
	"password-danie/internal/repository"
	"password-danie/internal/security"
)

type Vault struct {
	secrets repository.SecretRepo
}

func NewVault(secrets repository.SecretRepo) *Vault { return &Vault{secrets: secrets} }

func (v *Vault) Create(userID int64, username, passwordPlain, url, notes, icon string, title *string) (int64, error) {
	if username == "" {
		return 0, errors.New("username required")
	}
	cipher, iv, err := security.Encrypt([]byte(passwordPlain))
	if err != nil {
		return 0, err
	}
	t := ""
	if title != nil {
		t = *title
	}
	s := &domain.Secret{
		UserID:         userID,
		Username:       username,
		URL:            url,
		URLDomain:      extractDomain(url),
		Notes:          notes,
		Icon:           icon,
		Title:          t,
		PasswordCipher: cipher,
		PasswordIV:     iv,
	}
	return v.secrets.Create(s)
}

func (v *Vault) Get(userID, id int64) (*domain.Secret, error) {
	return v.secrets.GetByID(userID, id)
}

func (v *Vault) List(userID int64, q, domain string, limit, offset int) ([]domain.Secret, int, error) {
	filter := repository.ListFilter{Q: q, Domain: domain, Limit: limit, Offset: offset}
	return v.secrets.List(userID, filter)
}

func (v *Vault) Update(userID, id int64, username, passwordPlain, url, notes, icon, title *string) error {
	// Fetch, mutate, then persist
	cur, err := v.secrets.GetByID(userID, id)
	if err != nil {
		return err
	}
	if cur == nil {
		return errors.New("not found")
	}
	if username != nil {
		cur.Username = *username
	}
	if url != nil {
		cur.URL = *url
		cur.URLDomain = extractDomain(cur.URL)
	}
	if notes != nil {
		cur.Notes = *notes
	}
	if icon != nil {
		cur.Icon = *icon
	}
	if title != nil {
		cur.Title = *title
	}
	if passwordPlain != nil {
		c, iv, err := security.Encrypt([]byte(*passwordPlain))
		if err != nil {
			return err
		}
		cur.PasswordCipher = c
		cur.PasswordIV = iv
	}
	return v.secrets.Update(cur)
}

func (v *Vault) Delete(userID, id int64) error {
	return v.secrets.Delete(userID, id)
}
