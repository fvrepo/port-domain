package storage

import (
	"context"

	"github.com/port-domain/internal/models"
)

func (s *Storage) InsertOrUpdatePort(ctx context.Context, port *models.Port) error {
	return nil
}

func (s *Storage) GetPort(ctx context.Context, id string) (*models.Port, error) {
	return &models.Port{}, nil
}
