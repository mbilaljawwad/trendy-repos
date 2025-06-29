package datastore

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mbilaljawwad/trendy-repos/internal/migration"
	"github.com/spf13/viper"
)

func NewDataStore(ctx context.Context) *sqlx.DB {
	db := NewDataStoreWithoutMigrations(ctx)

	// Run database migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Error running database migrations: %v", err)
	}

	return db
}

func NewDataStoreWithoutMigrations(ctx context.Context) *sqlx.DB {
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("DB_NAME")

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable client_encoding=UTF8", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sqlx.ConnectContext(ctx, "postgres", connString)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Connected to database")
	return db
}

func runMigrations(db *sqlx.DB) error {
	migrator := migration.NewMigrator(db, "migrations")

	log.Println("Running database migrations...")
	if err := migrator.Up(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
