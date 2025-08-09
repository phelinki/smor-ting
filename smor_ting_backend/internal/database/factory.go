package database

import (
	"context"
	"fmt"

	"github.com/smorting/backend/configs"
	"github.com/smorting/backend/pkg/database"
	"github.com/smorting/backend/pkg/logger"
	"go.uber.org/zap"
)

// NewRepository creates a new repository based on configuration
func NewRepository(config *configs.DatabaseConfig, logger *logger.Logger) (Repository, error) {
	if config.InMemory {
		logger.Info("Using in-memory database")
		return NewMemoryDatabase(), nil
	}

	// Use MongoDB
	mongoDB, err := database.NewMongoDB(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB connection: %w", err)
	}

	repo := NewMongoDBRepository(mongoDB.GetDB(), logger)

	// Setup indexes
	ctx := context.Background()
	if err := repo.SetupIndexes(ctx); err != nil {
		logger.Warn("Failed to setup indexes", zap.Error(err))
	}

	return repo, nil
}
