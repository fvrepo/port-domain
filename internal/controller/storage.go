package controller

import (
	"context"

	"github.com/port-domain/internal/models"
)

type Storage interface {
	InsertOrUpdatePort(ctx context.Context, port *models.Port) error
	GetPort(ctx context.Context, id string) (*models.Port, error)
}
