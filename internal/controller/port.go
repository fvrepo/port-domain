package controller

import (
	"context"

	"github.com/port-domain/internal/models"
)

func (c *Controller) SavePort(ctx context.Context, p *models.Port) error {
	if err := c.Storage.InsertOrUpdatePort(ctx, p); err != nil {
		return nil
	}
	return nil
}

func (c *Controller) GetAllPorts(ctx context.Context, limit, skip int) ([]*models.Port, error) {
	return c.Storage.GetPorts(ctx, limit, skip)
}
