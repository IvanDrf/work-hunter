package fixtures

import (
	"time"

	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/internal/domain/rules"
	"github.com/google/uuid"
)

func CreateUsers() []*models.User {
	return []*models.User{
		{
			ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			Email:       "john.doe@example.com",
			FirstName:   "John",
			LastName:    "Doe",
			CompanyName: "Tech Corp",
			Status:      rules.UserStatusActive,
			Role:        rules.UserRoleEmployee,
			Verificated: true,
			CreatedAt:   time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:          uuid.MustParse("660e8400-e29b-41d4-a716-446655440001"),
			Email:       "jane.smith@example.com",
			FirstName:   "Jane",
			LastName:    "Smith",
			CompanyName: "Dev Inc",
			Status:      rules.UserStatusActive,
			Role:        rules.UserRoleAdmin,
			Verificated: false,
			CreatedAt:   time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:          uuid.MustParse("770e8400-e29b-41d4-a716-446655440002"),
			Email:       "bob.wilson@example.com",
			FirstName:   "Bob",
			LastName:    "Wilson",
			CompanyName: "Startup Ltd",
			Status:      rules.UserStatusBlocked,
			Role:        rules.UserRoleEmployee,
			Verificated: true,
			CreatedAt:   time.Date(2024, 1, 3, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 3, 10, 0, 0, 0, time.UTC),
		},
	}
}

func CreateSingleUser() *models.User {
	return CreateUsers()[0]
}

func ListUsersParams() *models.ListUsersParams {
	return &models.ListUsersParams{
		Page:     1,
		PageSize: 10,
	}
}

func ListUsersParamsWithFilters() *models.ListUsersParams {
	return &models.ListUsersParams{
		Page:        1,
		PageSize:    10,
		Status:      string(rules.UserStatusActive),
		Role:        string(rules.UserRoleEmployee),
		SearchQuery: "john",
		SortBy:      "created_at",
		SortDesc:    true,
	}
}
