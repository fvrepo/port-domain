package controller

import (
	"context"

	"github.com/port-domain/internal/models"
)

type Storage interface {
	InsertOrUpdatePort(ctx context.Context, port *models.Port) error
	GetPorts(ctx context.Context, limit, skip int) ([]*models.Port, error)
}
