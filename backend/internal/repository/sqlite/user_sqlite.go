// Adaptador SQLite de UserRepo: operaciones básicas sobre la tabla users.
package sqlite

import (
	"database/sql"
	"errors"

	"password-danie/internal/domain"
	"password-danie/internal/repository"
)

type UserSQLite struct{ db *sql.DB }

func NewUserSQLite(db *sql.DB) repository.UserRepo { return &UserSQLite{db: db} }

func (r *UserSQLite) Create(email, passwordHash string) (int64, error) {
	res, err := r.db.Exec(`INSERT INTO users(email, password_hash) VALUES(?, ?)`, email, passwordHash)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *UserSQLite) GetByID(id int64) (*domain.User, error) {
	row := r.db.QueryRow(`SELECT id, email, password_hash, created_at, updated_at FROM users WHERE id = ?`, id)
	var u domain.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserSQLite) GetByEmail(email string) (*domain.User, error) {
	row := r.db.QueryRow(`SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = ?`, email)
	var u domain.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
