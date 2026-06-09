package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(ctx context.Context, db *pgxpool.Pool, migrationsDir string) error {
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for _, file := range files {
		var applied bool
		err := db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE filename = $1)`, file).Scan(&applied)
		if err != nil {
			return fmt.Errorf("check migration %s: %w", file, err)
		}
		if applied {
			continue
		}

		content, err := os.ReadFile(filepath.Join(migrationsDir, file))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", file, err)
		}

		sql := string(content)
		if _, err := db.Exec(ctx, sql); err != nil {
			return fmt.Errorf("apply migration %s: %w", file, err)
		}

		if _, err := db.Exec(ctx, `INSERT INTO schema_migrations (filename) VALUES ($1)`, file); err != nil {
			return fmt.Errorf("record migration %s: %w", file, err)
		}

		log.Printf("migration applied: %s", file)
	}

	return nil
}
