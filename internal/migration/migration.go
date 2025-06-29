package migration

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type Migration struct {
	Version   int64
	Name      string
	UpSQL     string
	DownSQL   string
	Timestamp time.Time
}

type Migrator struct {
	db             *sqlx.DB
	migrationsPath string
	tableName      string
}

// NewMigrator creates a new migration runner
func NewMigrator(db *sqlx.DB, migrationsPath string) *Migrator {
	return &Migrator{
		db:             db,
		migrationsPath: migrationsPath,
		tableName:      "schema_migrations",
	}
}

// CreateMigrationsTable creates the table to track applied migrations
func (m *Migrator) CreateMigrationsTable() error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			version BIGINT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`, m.tableName)

	_, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	log.Printf("Migrations table '%s' created or already exists", m.tableName)
	return nil
}

// GetAppliedMigrations returns a list of applied migration versions
func (m *Migrator) GetAppliedMigrations() (map[int64]bool, error) {
	query := fmt.Sprintf("SELECT version FROM %s", m.tableName)
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[int64]bool)
	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("failed to scan migration version: %w", err)
		}
		applied[version] = true
	}

	return applied, nil
}

// LoadMigrations loads all migration files from the migrations directory
func (m *Migrator) LoadMigrations() ([]Migration, error) {
	var migrations []Migration

	err := filepath.WalkDir(m.migrationsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}

		migration, err := m.parseMigrationFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse migration file %s: %w", path, err)
		}

		migrations = append(migrations, migration)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load migrations: %w", err)
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// parseMigrationFile parses a migration SQL file
func (m *Migrator) parseMigrationFile(filePath string) (Migration, error) {
	filename := filepath.Base(filePath)

	// Parse version from filename (format: YYYYMMDDHHMMSS_migration_name.sql)
	re := regexp.MustCompile(`^(\d{14})_(.+)\.sql$`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) != 3 {
		return Migration{}, fmt.Errorf("invalid migration filename format: %s", filename)
	}

	version, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return Migration{}, fmt.Errorf("invalid version in filename: %s", matches[1])
	}

	name := matches[2]

	content, err := os.ReadFile(filePath)
	if err != nil {
		return Migration{}, fmt.Errorf("failed to read migration file: %w", err)
	}

	// Split content into UP and DOWN sections
	sections := strings.Split(string(content), "-- +migrate Down")
	if len(sections) != 2 {
		return Migration{}, fmt.Errorf("migration file must contain both UP and DOWN sections")
	}

	upSQL := strings.TrimPrefix(sections[0], "-- +migrate Up\n")
	downSQL := strings.TrimSpace(sections[1])

	return Migration{
		Version: version,
		Name:    name,
		UpSQL:   strings.TrimSpace(upSQL),
		DownSQL: downSQL,
	}, nil
}

// Up runs all pending migrations
func (m *Migrator) Up() error {
	log.Println("Starting database migrations...")

	if err := m.CreateMigrationsTable(); err != nil {
		return err
	}

	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	migrations, err := m.LoadMigrations()
	if err != nil {
		return err
	}

	count := 0
	for _, migration := range migrations {
		if applied[migration.Version] {
			continue
		}

		log.Printf("Applying migration %d: %s", migration.Version, migration.Name)

		tx, err := m.db.Beginx()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		// Execute the migration
		if _, err := tx.Exec(migration.UpSQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %d: %w", migration.Version, err)
		}

		// Record the migration as applied
		query := fmt.Sprintf("INSERT INTO %s (version, name) VALUES ($1, $2)", m.tableName)
		if _, err := tx.Exec(query, migration.Version, migration.Name); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
		}

		log.Printf("Successfully applied migration %d: %s", migration.Version, migration.Name)
		count++
	}

	if count == 0 {
		log.Println("No pending migrations to apply")
	} else {
		log.Printf("Successfully applied %d migrations", count)
	}

	return nil
}

// Down rolls back the last applied migration
func (m *Migrator) Down() error {
	log.Println("Rolling back last migration...")

	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	if len(applied) == 0 {
		log.Println("No migrations to roll back")
		return nil
	}

	migrations, err := m.LoadMigrations()
	if err != nil {
		return err
	}

	// Find the last applied migration
	var lastMigration *Migration
	for i := len(migrations) - 1; i >= 0; i-- {
		if applied[migrations[i].Version] {
			lastMigration = &migrations[i]
			break
		}
	}

	if lastMigration == nil {
		log.Println("No migrations to roll back")
		return nil
	}

	log.Printf("Rolling back migration %d: %s", lastMigration.Version, lastMigration.Name)

	tx, err := m.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Execute the rollback
	if _, err := tx.Exec(lastMigration.DownSQL); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute rollback %d: %w", lastMigration.Version, err)
	}

	// Remove the migration record
	query := fmt.Sprintf("DELETE FROM %s WHERE version = $1", m.tableName)
	if _, err := tx.Exec(query, lastMigration.Version); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record %d: %w", lastMigration.Version, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit rollback %d: %w", lastMigration.Version, err)
	}

	log.Printf("Successfully rolled back migration %d: %s", lastMigration.Version, lastMigration.Name)
	return nil
}

// Status shows the current migration status
func (m *Migrator) Status() error {
	log.Println("Migration status:")

	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	migrations, err := m.LoadMigrations()
	if err != nil {
		return err
	}

	if len(migrations) == 0 {
		log.Println("No migration files found")
		return nil
	}

	for _, migration := range migrations {
		status := "Pending"
		if applied[migration.Version] {
			status = "Applied"
		}
		log.Printf("  %d: %s [%s]", migration.Version, migration.Name, status)
	}

	return nil
}

// CreateMigrationFile creates a new migration file with the given name
func (m *Migrator) CreateMigrationFile(name string) error {
	if err := os.MkdirAll(m.migrationsPath, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s.sql", timestamp, name)
	filepath := filepath.Join(m.migrationsPath, filename)

	template := `-- +migrate Up
-- Write your UP migration here

-- +migrate Down
-- Write your DOWN migration here
`

	if err := os.WriteFile(filepath, []byte(template), 0644); err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}

	log.Printf("Created migration file: %s", filepath)
	return nil
}
