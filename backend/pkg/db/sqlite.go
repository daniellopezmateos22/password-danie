package db

// Package db ofrece utilidades para abrir y configurar la conexión SQLite.

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func OpenSQLite(dsn string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", dsn+"?_foreign_keys=on&_busy_timeout=5000")
    if err != nil {
        return nil, err
    }
    if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
        return nil, err
    }
    return db, nil
}
