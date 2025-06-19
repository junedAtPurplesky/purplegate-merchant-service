package db

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

// Migration represents a database migration
type Migration struct {
	Version int
	UpSQL   string
	DownSQL string
}

// RunMigrations runs database migrations
func RunMigrations(pool *pgxpool.Pool, migrationsPath string) error {
	log.Println("Running database migrations...")

	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(pool); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	migrations, err := loadMigrations(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := getAppliedMigrations(pool)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Apply pending migrations
	for _, migration := range migrations {
		if !appliedMigrations[migration.Version] {
			log.Printf("Applying migration %d", migration.Version)
			if err := applyMigration(pool, migration); err != nil {
				return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
			}
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// createMigrationsTable creates the migrations tracking table
func createMigrationsTable(pool *pgxpool.Pool) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`
	_, err := pool.Exec(context.Background(), query)
	return err
}

// loadMigrations loads migration files from the migrations directory
func loadMigrations(migrationsPath string) ([]Migration, error) {
	var migrations []Migration

	// Read migration files
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".up.sql") {
			continue
		}

		// Parse version from filename (e.g., "000001_create_merchants_table.up.sql")
		parts := strings.Split(file.Name(), "_")
		if len(parts) < 2 {
			continue
		}

		var version int
		if _, err := fmt.Sscanf(parts[0], "%d", &version); err != nil {
			continue
		}

		// Read up migration
		upPath := filepath.Join(migrationsPath, file.Name())
		upSQL, err := os.ReadFile(upPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", upPath, err)
		}

		// Read down migration
		downPath := strings.Replace(upPath, ".up.sql", ".down.sql", 1)
		downSQL, err := os.ReadFile(downPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read down migration file %s: %w", downPath, err)
		}

		migrations = append(migrations, Migration{
			Version: version,
			UpSQL:   string(upSQL),
			DownSQL: string(downSQL),
		})
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// getAppliedMigrations gets the list of already applied migrations
func getAppliedMigrations(pool *pgxpool.Pool) (map[int]bool, error) {
	query := `SELECT version FROM schema_migrations ORDER BY version`
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// applyMigration applies a single migration
func applyMigration(pool *pgxpool.Pool, migration Migration) error {
	ctx := context.Background()

	// Start transaction
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Execute migration SQL
	if _, err := tx.Exec(ctx, migration.UpSQL); err != nil {
		return err
	}

	// Record migration as applied
	recordQuery := `INSERT INTO schema_migrations (version) VALUES ($1)`
	if _, err := tx.Exec(ctx, recordQuery, migration.Version); err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit(ctx)
}
