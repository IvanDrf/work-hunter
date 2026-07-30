package user

import (
	"testing"

	"github.com/IvanDrf/work-hunter/users/internal/config"
	"github.com/IvanDrf/work-hunter/users/internal/infrastructure/service"
	"github.com/IvanDrf/work-hunter/users/internal/logger"
	"github.com/IvanDrf/work-hunter/users/tests/mocks"
)

func setupTest(t *testing.T) (*service.UserService, *mocks.UserRepositoryMock) {
	mockRepo := new(mocks.UserRepositoryMock)
	logger := logger.New(&config.LoggerConfig{Level: "debug"})
	userService := service.NewUserService(mockRepo, logger)
	return userService, mockRepo
}
