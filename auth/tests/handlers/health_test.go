package handlers_test

import (
	"testing"

	"github.com/IvanDrf/work-hunter/pkg/common"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	t.Parallel()

	handlers := newHandlers(nil)

	resp, err := handlers.Health(t.Context(), &common.Empty{})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.Status, common.Status_AVAILABLE)
}
