package controller

import (
	"context"
	"testing"

	"github.com/port-domain/internal/models"

	"github.com/stretchr/testify/require"
)

func TestStartContainer(t *testing.T) {
	ctx := context.Background()
	mg, err := NewDockerMongoStorage()
	require.NoError(t, err)
	require.NotNil(t, mg)

	p1 := &models.Port{
		ID:   "test",
		Code: "123123",
		City: "kiev",
	}

	p2 := &models.Port{
		ID:      "test",
		Code:    "123123",
		City:    "kiev",
		Country: "ukraine",
	}

	err = mg.InsertOrUpdatePort(ctx, p1)
	require.NoError(t, err)

	var ps []*models.Port
	ps, err = mg.GetPorts(ctx, 10)
	require.NoError(t, err)
	require.NotNil(t, ps)
	require.Equal(t, 1, len(ps))
	require.Equal(t, p1.ID, ps[0].ID)
	require.Equal(t, p1.Code, ps[0].Code)
	require.Equal(t, p1.City, ps[0].City)
	require.Equal(t, p1.Country, ps[0].Country)

	err = mg.InsertOrUpdatePort(ctx, p2)
	require.NoError(t, err)

	ps, err = mg.GetPorts(ctx, 10)
	require.NoError(t, err)
	require.NotNil(t, ps)
	require.Equal(t, 1, len(ps))
	require.Equal(t, p2.ID, ps[0].ID)
	require.Equal(t, p2.Code, ps[0].Code)
	require.Equal(t, p2.City, ps[0].City)
	require.Equal(t, p2.Country, ps[0].Country)

}
