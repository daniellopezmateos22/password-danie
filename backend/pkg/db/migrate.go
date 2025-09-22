// Package db: utilidades de migración sencilla leyendo .sql de un directorio.
package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ApplyMigrations(sqlDB *sql.DB, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		// si no hay carpeta de migraciones, no falles en dev
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	// ordenar por nombre para aplicar 001_, 002_, ...
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", path, err)
		}
		sqlText := string(b)
		stmts := splitSQL(sqlText)
		for _, s := range stmts {
			if strings.TrimSpace(s) == "" {
				continue
			}
			if _, err := sqlDB.Exec(s); err != nil {
				return fmt.Errorf("exec migration %s: %w", path, err)
			}
		}
	}
	return nil
}

func splitSQL(sql string) []string {
	parts := strings.Split(sql, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
