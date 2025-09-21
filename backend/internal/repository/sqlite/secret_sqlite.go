// Adaptador SQLite de SecretRepo: CRUD y listado con búsqueda/filtro de dominio.
package sqlite

import (
	"database/sql"
	"fmt"
	"strings"

	"password-danie/internal/domain"
	"password-danie/internal/repository"
)

type SecretSQLite struct{ db *sql.DB }

func NewSecretSQLite(db *sql.DB) repository.SecretRepo { return &SecretSQLite{db: db} }

func (r *SecretSQLite) Create(s *domain.Secret) (int64, error) {
	res, err := r.db.Exec(`INSERT INTO secrets(user_id, username, password_cipher, password_iv, url, url_domain, notes, icon, title)
                           VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		s.UserID, s.Username, s.PasswordCipher, s.PasswordIV, s.URL, s.URLDomain, s.Notes, s.Icon, s.Title)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *SecretSQLite) GetByID(userID, id int64) (*domain.Secret, error) {
	row := r.db.QueryRow(`SELECT id, user_id, username, password_cipher, password_iv, url, url_domain, notes, icon, title, created_at, updated_at
	                      FROM secrets WHERE id = ? AND user_id = ?`, id, userID)
	var s domain.Secret
	if err := row.Scan(&s.ID, &s.UserID, &s.Username, &s.PasswordCipher, &s.PasswordIV, &s.URL, &s.URLDomain, &s.Notes, &s.Icon, &s.Title, &s.CreatedAt, &s.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *SecretSQLite) List(userID int64, f repository.ListFilter) ([]domain.Secret, int, error) {
	where := []string{"user_id = ?"}
	args := []any{userID}

	q := strings.TrimSpace(f.Q)
	if q != "" {
		where = append(where, "(username LIKE ? OR url LIKE ? OR notes LIKE ? OR title LIKE ?)")
		like := "%" + q + "%"
		args = append(args, like, like, like, like)
	}
	if d := strings.TrimSpace(f.Domain); d != "" {
		where = append(where, "url_domain = ?")
		args = append(args, d)
	}
	if f.Limit <= 0 {
		f.Limit = 20
	}
	if f.Offset < 0 {
		f.Offset = 0
	}

	w := "WHERE " + strings.Join(where, " AND ")
	var total int
	if err := r.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM secrets %s", w), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(fmt.Sprintf(`SELECT id, user_id, username, password_cipher, password_iv, url, url_domain, notes, icon, title, created_at, updated_at
	                                    FROM secrets %s ORDER BY id DESC LIMIT ? OFFSET ?`, w),
		append(args, f.Limit, f.Offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []domain.Secret
	for rows.Next() {
		var s domain.Secret
		if err := rows.Scan(&s.ID, &s.UserID, &s.Username, &s.PasswordCipher, &s.PasswordIV, &s.URL, &s.URLDomain, &s.Notes, &s.Icon, &s.Title, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, s)
	}
	return out, total, rows.Err()
}

func (r *SecretSQLite) Update(s *domain.Secret) error {
	_, err := r.db.Exec(`UPDATE secrets
	                     SET username=?, password_cipher=?, password_iv=?, url=?, url_domain=?, notes=?, icon=?, title=?, updated_at=CURRENT_TIMESTAMP
	                     WHERE id=? AND user_id=?`,
		s.Username, s.PasswordCipher, s.PasswordIV, s.URL, s.URLDomain, s.Notes, s.Icon, s.Title, s.ID, s.UserID)
	return err
}

func (r *SecretSQLite) Delete(userID, id int64) error {
	_, err := r.db.Exec(`DELETE FROM secrets WHERE id=? AND user_id=?`, id, userID)
	return err
}
