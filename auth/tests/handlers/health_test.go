package handlers_test

import (
	"testing"

	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestHealth(t *testing.T) {
	t.Parallel()

	handlers := newHandlers(nil)

	resp, err := handlers.Health(t.Context(), &auth_api.Empty{})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.Code, int32(codes.OK))
}
