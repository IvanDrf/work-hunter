package handlers_test

import (
	"testing"

	"github.com/IvanDrf/work-hunter/pkg/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	t.Parallel()

	handlers := newHandlers(nil)

	resp, err := handlers.Health(t.Context(), &common.Empty{})
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, common.Status_AVAILABLE, resp.GetStatus())
}
