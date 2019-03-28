package controller

import (
	"context"

	"github.com/port-domain/internal/models"
)

func (c *Controller) SavePort(ctx context.Context, p *models.Port) error {
	return nil
}

func (c *Controller) GetAllPorts(ctx context.Context, limit int, cursor string) ([]*models.Port, error) {
	return nil, nil
}
