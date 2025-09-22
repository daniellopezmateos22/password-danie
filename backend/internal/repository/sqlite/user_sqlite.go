// Adaptador SQLite de UserRepo: operaciones básicas + soporte de password reset (con valores, no punteros).
package sqlite

import (
	"database/sql"
	"errors"
	"time"

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
	row := r.db.QueryRow(`SELECT id, email, password_hash, reset_token, reset_expires_at, created_at, updated_at FROM users WHERE id = ?`, id)
	var u domain.User
	var resetToken sql.NullString
	var resetExpires sql.NullTime

	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &resetToken, &resetExpires, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if resetToken.Valid {
		u.ResetToken = resetToken.String
	} else {
		u.ResetToken = ""
	}
	if resetExpires.Valid {
		u.ResetExpiresAt = resetExpires.Time
	} else {
		u.ResetExpiresAt = time.Time{}
	}
	return &u, nil
}

func (r *UserSQLite) GetByEmail(email string) (*domain.User, error) {
	row := r.db.QueryRow(`SELECT id, email, password_hash, reset_token, reset_expires_at, created_at, updated_at FROM users WHERE email = ?`, email)
	var u domain.User
	var resetToken sql.NullString
	var resetExpires sql.NullTime

	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &resetToken, &resetExpires, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if resetToken.Valid {
		u.ResetToken = resetToken.String
	} else {
		u.ResetToken = ""
	}
	if resetExpires.Valid {
		u.ResetExpiresAt = resetExpires.Time
	} else {
		u.ResetExpiresAt = time.Time{}
	}
	return &u, nil
}

// --- password reset ---

func (r *UserSQLite) UpdateReset(userID int64, token *string, expiresAt *time.Time) error {
	_, err := r.db.Exec(`UPDATE users SET reset_token = ?, reset_expires_at = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		token, expiresAt, userID)
	return err
}

func (r *UserSQLite) GetByResetToken(token string) (*domain.User, error) {
	row := r.db.QueryRow(`SELECT id, email, password_hash, reset_token, reset_expires_at, created_at, updated_at FROM users WHERE reset_token = ?`, token)
	var u domain.User
	var resetToken sql.NullString
	var resetExpires sql.NullTime

	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &resetToken, &resetExpires, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if resetToken.Valid {
		u.ResetToken = resetToken.String
	} else {
		u.ResetToken = ""
	}
	if resetExpires.Valid {
		u.ResetExpiresAt = resetExpires.Time
	} else {
		u.ResetExpiresAt = time.Time{}
	}
	return &u, nil
}

func (r *UserSQLite) UpdatePassword(userID int64, passwordHash string) error {
	_, err := r.db.Exec(`UPDATE users 
		SET password_hash = ?, reset_token = NULL, reset_expires_at = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		passwordHash, userID)
	return err
}
