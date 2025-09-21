// Package repository declara puertos (interfaces) para secretos del vault y filtros de listado.
package repository

import "password-danie/internal/domain"

type ListFilter struct {
	Q      string
	Domain string 
	Limit  int
	Offset int
}

type SecretRepo interface {
	Create(s *domain.Secret) (int64, error)
	GetByID(userID, id int64) (*domain.Secret, error)
	List(userID int64, f ListFilter) ([]domain.Secret, int, error)
	Update(s *domain.Secret) error
	Delete(userID, id int64) error
}
