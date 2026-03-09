package repo

import "github.com/IvanDrf/work-hunter/auth/internal/domain/models"

type UserRepo interface {
	CreateUser(user *models.User) error
	FindUser(username string) (*models.User, error)
}
