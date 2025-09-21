// Caso de uso de autenticación: registro (bcrypt) y login (bcrypt + JWT).
package usecase

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"password-danie/internal/domain"
	"password-danie/internal/repository"
	"password-danie/internal/security"
)

type Auth struct {
	users repository.UserRepo
}

func NewAuth(users repository.UserRepo) *Auth { return &Auth{users: users} }

func (a *Auth) Register(email, password string) (*domain.User, error) {
	if len(password) < 8 {
		return nil, errors.New("password too short")
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	id, err := a.users.Create(email, string(hash))
	if err != nil {
		return nil, err
	}
	return a.users.GetByID(id)
}

func (a *Auth) Login(email, password string) (string, *domain.User, error) {
	u, err := a.users.GetByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if u == nil {
		return "", nil, errors.New("invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", nil, errors.New("invalid credentials")
	}
	token, err := security.GenerateAccessToken(u.ID, u.Email, 15*time.Minute)
	return token, u, err
}
