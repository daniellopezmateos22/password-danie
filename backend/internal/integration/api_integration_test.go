// Test de integración end-to-end: recorre todos los endpoints (auth, users/me, vault CRUD, reset).
package integration_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/gin-gonic/gin"

	api "password-danie/internal/http"
	"password-danie/internal/repository"
	sqlrepo "password-danie/internal/repository/sqlite"
	"password-danie/internal/usecase"
)

// ---------- helpers ----------

func mustStatus(t *testing.T, rr *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rr.Code != want {
		t.Fatalf("status %d, want %d. body=%s", rr.Code, want, rr.Body.String())
	}
}

func doJSON(t *testing.T, ts *httptest.Server, method, path string, token string, payload any) *httptest.ResponseRecorder {
	t.Helper()
	var body []byte
	var err error
	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
	}
	req, _ := http.NewRequest(method, ts.URL+path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rr := httptest.NewRecorder()
	ts.Config.Handler.ServeHTTP(rr, req)
	return rr
}

type regRes struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}
type loginRes struct {
	AccessToken string `json:"access_token"`
	User        struct {
		ID    int64  `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
}
type createRes struct {
	ID int64 `json:"id"`
}
type listRes struct {
	Items []struct {
		ID int64 `json:"id"`
	} `json:"items"`
	Total int `json:"total"`
}

// ---------- schema (mimic migrations) ----------
const schema = `
CREATE TABLE IF NOT EXISTS users(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  reset_token TEXT NULL,
  reset_expires_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS secrets(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  username TEXT NOT NULL,
  password_cipher TEXT NOT NULL,
  password_iv TEXT NOT NULL,
  url TEXT NOT NULL DEFAULT '',
  url_domain TEXT NOT NULL DEFAULT '',
  notes TEXT NOT NULL DEFAULT '',
  icon TEXT NOT NULL DEFAULT '',
  title TEXT NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

// ---------- test ----------

func Test_FullAPI_HappyPath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Env necesarios para JWT y AES
	_ = os.Setenv("JWT_SECRET", "test-secret")
	_ = os.Setenv("AES_KEY", "0123456789abcdef0123456789abcdef")

	// DB en memoria
	sqlDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite mem: %v", err)
	}
	defer sqlDB.Close()
	if _, err := sqlDB.Exec(schema); err != nil {
		t.Fatalf("apply schema: %v", err)
	}

	// Repos & Usecases
	var (
		userRepo   repository.UserRepo   = sqlrepo.NewUserSQLite(sqlDB)
		secretRepo repository.SecretRepo = sqlrepo.NewSecretSQLite(sqlDB)
	)
	authUC := usecase.NewAuth(userRepo)
	vaultUC := usecase.NewVault(secretRepo)
	resetUC := usecase.NewPasswordReset(userRepo)

	// Router y server
	r := gin.Default()
	api.RegisterRoutes(r, authUC, vaultUC, func() error { return sqlDB.Ping() })
	api.RegisterResetRoutes(r, resetUC)
	ts := httptest.NewServer(r)
	defer ts.Close()

	// --- 1) health/ready
	rr := doJSON(t, ts, http.MethodGet, "/healthz", "", nil)
	mustStatus(t, rr, 200)
	rr = doJSON(t, ts, http.MethodGet, "/readyz", "", nil)
	mustStatus(t, rr, 200)

	// --- 2) register
	registerBody := map[string]any{"email": "e2e@test.com", "password": "Secret123!"}
	rr = doJSON(t, ts, http.MethodPost, "/api/v1/auth/register", "", registerBody)
	mustStatus(t, rr, 201)
	var reg regRes
	_ = json.Unmarshal(rr.Body.Bytes(), &reg)
	if reg.ID == 0 || reg.Email != "e2e@test.com" {
		t.Fatalf("bad register response: %+v", reg)
	}

	// --- 3) login
	rr = doJSON(t, ts, http.MethodPost, "/api/v1/auth/login", "", registerBody)
	mustStatus(t, rr, 200)
	var logres loginRes
	_ = json.Unmarshal(rr.Body.Bytes(), &logres)
	if logres.AccessToken == "" || logres.User.ID == 0 {
		t.Fatalf("bad login response: %s", rr.Body.String())
	}
	token := logres.AccessToken

	// --- 4) users/me
	rr = doJSON(t, ts, http.MethodGet, "/api/v1/users/me", token, nil)
	mustStatus(t, rr, 200)

	// --- 5) crear secreto
	createBody := map[string]any{
		"username":       "danie",
		"password_plain": "p@ss",
		"url":            "https://github.com",
		"notes":          "mi cuenta",
		"icon":           "github",
		"title":          "GitHub",
	}
	rr = doJSON(t, ts, http.MethodPost, "/api/v1/vault/entries", token, createBody)
	mustStatus(t, rr, 201)
	var cres createRes
	_ = json.Unmarshal(rr.Body.Bytes(), &cres)
	if cres.ID == 0 {
		t.Fatalf("no id returned: %s", rr.Body.String())
	}

	// --- 6) listar
	rr = doJSON(t, ts, http.MethodGet, "/api/v1/vault/entries?limit=20&offset=0&q=git", token, nil)
	mustStatus(t, rr, 200)
	var lres listRes
	_ = json.Unmarshal(rr.Body.Bytes(), &lres)
	if lres.Total == 0 || len(lres.Items) == 0 {
		t.Fatalf("expected at least 1 item: %s", rr.Body.String())
	}

	// --- 7) obtener por id
	idStr := strconv.FormatInt(cres.ID, 10)
	rr = doJSON(t, ts, http.MethodGet, "/api/v1/vault/entries/"+idStr, token, nil)
	mustStatus(t, rr, 200)

	// --- 8) actualizar
	updateBody := map[string]any{
		"notes":          "actualizada",
		"password_plain": "new-pass-123",
	}
	rr = doJSON(t, ts, http.MethodPut, "/api/v1/vault/entries/"+idStr, token, updateBody)
	mustStatus(t, rr, 200)

	// --- 9) borrar
	rr = doJSON(t, ts, http.MethodDelete, "/api/v1/vault/entries/"+idStr, token, nil)
	mustStatus(t, rr, 200)

	// --- 10) reset password: request
	rr = doJSON(t, ts, http.MethodPost, "/api/v1/auth/reset/request", "", map[string]any{"email": "e2e@test.com"})
	mustStatus(t, rr, 200)
	var rres map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &rres)
	resetToken := rres["reset_token"]
	if resetToken == "" {
		t.Fatalf("no reset_token: %s", rr.Body.String())
	}

	// --- 11) reset password: confirm
	rr = doJSON(t, ts, http.MethodPost, "/api/v1/auth/reset/confirm", "", map[string]any{
		"token":        resetToken,
		"new_password": "NewPassw0rd!",
	})
	mustStatus(t, rr, 200)

	// --- 12) login con la nueva contraseña
	rr = doJSON(t, ts, http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"email":    "e2e@test.com",
		"password": "NewPassw0rd!",
	})
	mustStatus(t, rr, 200)
}
