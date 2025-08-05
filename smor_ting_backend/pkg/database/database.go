package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/smorting/backend/configs"
	"github.com/smorting/backend/pkg/logger"
	"go.uber.org/zap"
)

// Database represents a database connection
type Database struct {
	DB     *sql.DB
	Config *configs.DatabaseConfig
	Logger *logger.Logger
}

// New creates a new database connection
func New(config *configs.DatabaseConfig, log *logger.Logger) (*Database, error) {
	db := &Database{
		Config: config,
		Logger: log,
	}

	if err := db.connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// connect establishes a connection to the database
func (d *Database) connect() error {
	var dsn string

	switch d.Config.Driver {
	case "sqlite", "sqlite3":
		if d.Config.InMemory {
			dsn = ":memory:"
			d.Logger.Info("Using in-memory SQLite database (development/testing mode)")
		} else {
			dsn = d.Config.Database
			d.Logger.Info("Using SQLite database", zap.String("database", d.Config.Database))
		}
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			d.Config.Host, d.Config.Port, d.Config.Username, d.Config.Password,
			d.Config.Database, d.Config.SSLMode)
		d.Logger.Info("Using PostgreSQL database", zap.String("host", d.Config.Host))
	default:
		return fmt.Errorf("unsupported database driver: %s", d.Config.Driver)
	}

	db, err := sql.Open(d.Config.Driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	d.DB = db
	return nil
}

// ping verifies the database connection
func (d *Database) ping() error {
	if err := d.DB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	d.Logger.Info("Database connection established successfully")
	return nil
}

// migrate runs database migrations
func (d *Database) migrate() error {
	// For now, we'll create basic tables
	// In a real application, you'd use a migration tool like golang-migrate
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			first_name TEXT,
			last_name TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS services (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			category TEXT,
			price DECIMAL(10,2),
			user_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
	}

	for i, migration := range migrations {
		if _, err := d.DB.Exec(migration); err != nil {
			return fmt.Errorf("failed to execute migration %d: %w", i+1, err)
		}
	}

	d.Logger.Info("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.DB != nil {
		if err := d.DB.Close(); err != nil {
			d.Logger.Error("Failed to close database connection", err)
			return err
		}
		d.Logger.Info("Database connection closed")
	}
	return nil
}

// GetDB returns the underlying sql.DB instance
func (d *Database) GetDB() *sql.DB {
	return d.DB
}

// IsInMemory returns true if the database is in-memory
func (d *Database) IsInMemory() bool {
	return d.Config.InMemory
}

// HealthCheck performs a health check on the database
func (d *Database) HealthCheck() error {
	return d.DB.Ping()
}
