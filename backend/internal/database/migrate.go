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

func RebuildUserScores(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO user_scores (user_id, total_points, predictions_count, exact_scores_count,
		                         winner_correct_count, goal_difference_correct_count,
		                         streak_count, last_scored_at, updated_at)
		SELECT
			u.id,
			COALESCE(p_stats.total_points, 0),
			COALESCE(p_stats.predictions_count, 0),
			COALESCE(p_stats.exact_scores_count, 0),
			COALESCE(p_stats.winner_correct_count, 0),
			COALESCE(p_stats.goal_difference_correct_count, 0),
			0,
			se.last_scored_at,
			NOW()
		FROM users u
		LEFT JOIN (
			SELECT user_id,
			       SUM(total_points) AS total_points,
			       COUNT(id) FILTER (WHERE total_points > 0) AS predictions_count,
			       COUNT(id) FILTER (WHERE is_exact_score = true) AS exact_scores_count,
			       COUNT(id) FILTER (WHERE is_winner_correct = true) AS winner_correct_count,
			       COUNT(id) FILTER (WHERE is_goal_difference_correct = true) AS goal_difference_correct_count
			FROM predictions
			GROUP BY user_id
		) p_stats ON p_stats.user_id = u.id
		LEFT JOIN (
			SELECT user_id, MAX(created_at) AS last_scored_at
			FROM score_events
			GROUP BY user_id
		) se ON se.user_id = u.id
		ON CONFLICT (user_id) DO UPDATE SET
			total_points = EXCLUDED.total_points,
			predictions_count = EXCLUDED.predictions_count,
			exact_scores_count = EXCLUDED.exact_scores_count,
			winner_correct_count = EXCLUDED.winner_correct_count,
			goal_difference_correct_count = EXCLUDED.goal_difference_correct_count,
			last_scored_at = EXCLUDED.last_scored_at,
			updated_at = NOW()
	`
	_, err := db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("rebuild user scores: %w", err)
	}
	log.Println("user scores rebuilt successfully")
	return nil
}

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
