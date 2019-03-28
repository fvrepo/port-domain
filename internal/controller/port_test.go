package controller

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStartContainer(t *testing.T) {
	mg, err := NewDockerMongoStorage()
	require.NoError(t, err)
	require.NotNil(t, mg)
}
