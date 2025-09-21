-- Migración inicial: users + secrets, soporte de reset y filtrado por dominio.
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    reset_token TEXT NULL,
    reset_expires_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS secrets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    username TEXT NOT NULL,
    password_cipher TEXT NOT NULL,
    password_iv TEXT NOT NULL,     
    url TEXT NOT NULL DEFAULT '',
    url_domain TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    icon TEXT NOT NULL DEFAULT '',
    title TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_secrets_user    ON secrets(user_id);
CREATE INDEX IF NOT EXISTS idx_secrets_domain  ON secrets(url_domain);
CREATE INDEX IF NOT EXISTS idx_secrets_search  ON secrets(username, url, notes, title);
