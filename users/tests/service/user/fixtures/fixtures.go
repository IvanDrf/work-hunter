package fixtures

import (
	"time"

	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/internal/domain/rules"
	"github.com/IvanDrf/work-hunter/users/internal/interfaces/grpc/dto"
	"github.com/google/uuid"
)

func TestUser() *models.User {
	return &models.User{
		ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		Email:       "test@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		CompanyName: "Test Corp",
		Status:      rules.UserStatusActive,
		Role:        rules.UserRoleEmployer,
		Verificated: true,
		CreatedAt:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
	}
}

func TestUserEmployee() *models.User {
	return &models.User{
		ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		Email:       "employee@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		CompanyName: "",
		Status:      rules.UserStatusActive,
		Role:        rules.UserRoleEmployee,
		Verificated: true,
		CreatedAt:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
	}
}

func TestUserEmployer() *models.User {
	return &models.User{
		ID:          uuid.MustParse("660e8400-e29b-41d4-a716-446655440001"),
		Email:       "employer@example.com",
		FirstName:   "Jane",
		LastName:    "Smith",
		CompanyName: "Tech Corp",
		Status:      rules.UserStatusActive,
		Role:        rules.UserRoleEmployer,
		Verificated: false,
		CreatedAt:   time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC),
	}
}

func TestUser2() *models.User {
	return &models.User{
		ID:          uuid.MustParse("660e8400-e29b-41d4-a716-446655440001"),
		Email:       "jane@example.com",
		FirstName:   "Jane",
		LastName:    "Smith",
		CompanyName: "Test Inc",
		Status:      rules.UserStatusActive,
		Role:        rules.UserRoleEmployer,
		Verificated: false,
		CreatedAt:   time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC),
	}
}

func TestUser3() *models.User {
	return &models.User{
		ID:          uuid.MustParse("770e8400-e29b-41d4-a716-446655440002"),
		Email:       "bob@example.com",
		FirstName:   "Bob",
		LastName:    "Wilson",
		CompanyName: "Startup Ltd",
		Status:      rules.UserStatusBlocked,
		Role:        rules.UserRoleEmployer,
		Verificated: false,
		CreatedAt:   time.Date(2024, 1, 3, 12, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2024, 1, 3, 12, 0, 0, 0, time.UTC),
	}
}

func TestUsers() []*models.User {
	return []*models.User{
		TestUser(),
		TestUser2(),
		TestUser3(),
	}
}

func CreateUserRequestDTO() *dto.CreateUserRequest {
	return &dto.CreateUserRequest{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Email:       "test@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		CompanyName: "Test Corp",
		Verificated: true,
	}
}

func CreateUserRequestWithoutCompanyDTO() *dto.CreateUserRequest {
	return &dto.CreateUserRequest{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Email:       "employee@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		CompanyName: "",
		Verificated: true,
	}
}

func CreateUserRequestEmptyFieldsDTO() *dto.CreateUserRequest {
	return &dto.CreateUserRequest{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Email:       "test@example.com",
		FirstName:   "",
		LastName:    "",
		CompanyName: "",
		Verificated: false,
	}
}

func UpdateUserRequestDTO() *dto.UpdateUserRequest {
	firstName := "Updated"
	lastName := "Name"
	companyName := "New Corp"
	verificated := false

	return &dto.UpdateUserRequest{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		FirstName:   firstName,
		LastName:    lastName,
		CompanyName: companyName,
		Verificated: verificated,
	}
}

func UpdateUserRequestPartialDTO() *dto.UpdateUserRequest {
	firstName := "NewFirstName"

	return &dto.UpdateUserRequest{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		FirstName:   firstName,
		LastName:    "",
		CompanyName: "",
		Verificated: false,
	}
}

func UpdateUserRequestOnlyCompanyDTO() *dto.UpdateUserRequest {
	companyName := "Updated Corp"

	return &dto.UpdateUserRequest{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		FirstName:   "",
		LastName:    "",
		CompanyName: companyName,
		Verificated: false,
	}
}

func UpdateUserRequestInvalidDTO() *dto.UpdateUserRequest {
	return &dto.UpdateUserRequest{
		ID: "invalid-uuid",
	}
}

func ListUsersRequestDTO() *dto.ListUsersRequest {
	return &dto.ListUsersRequest{
		Page:     1,
		PageSize: 10,
		Filter: &dto.UserFilter{
			Status:      string(rules.UserStatusActive),
			Role:        string(rules.UserRoleEmployer),
			SearchQuery: "john",
		},
		SortBy:   "created_at",
		SortDesc: true,
	}
}

func ListUsersRequestWithoutFilterDTO() *dto.ListUsersRequest {
	return &dto.ListUsersRequest{
		Page:     1,
		PageSize: 10,
		Filter:   nil,
	}
}

func ListUsersRequestWithSearchDTO() *dto.ListUsersRequest {
	return &dto.ListUsersRequest{
		Page:     1,
		PageSize: 10,
		Filter: &dto.UserFilter{
			SearchQuery: "john",
		},
	}
}

func ListUsersRequestWithStatusDTO() *dto.ListUsersRequest {
	return &dto.ListUsersRequest{
		Page:     1,
		PageSize: 10,
		Filter: &dto.UserFilter{
			Status: string(rules.UserStatusActive),
		},
	}
}

func ListUsersRequestWithRoleDTO() *dto.ListUsersRequest {
	return &dto.ListUsersRequest{
		Page:     1,
		PageSize: 10,
		Filter: &dto.UserFilter{
			Role: string(rules.UserRoleEmployer),
		},
	}
}

func ListUsersRequestInvalidPaginationDTO() *dto.ListUsersRequest {
	return &dto.ListUsersRequest{
		Page:     0,
		PageSize: 0,
	}
}

func ListUsersRequestMaxPageSizeDTO() *dto.ListUsersRequest {
	return &dto.ListUsersRequest{
		Page:     1,
		PageSize: 150,
	}
}

func UpdateUserStatusRequestDTO() *dto.UpdateUserStatusRequest {
	return &dto.UpdateUserStatusRequest{
		ID:     "550e8400-e29b-41d4-a716-446655440000",
		Status: string(rules.UserStatusBlocked),
	}
}

func UpdateUserStatusRequestInvalidDTO() *dto.UpdateUserStatusRequest {
	return &dto.UpdateUserStatusRequest{
		ID:     "invalid-uuid",
		Status: string(rules.UserStatusBlocked),
	}
}

func UpdateUserStatusRequestActiveDTO() *dto.UpdateUserStatusRequest {
	return &dto.UpdateUserStatusRequest{
		ID:     "550e8400-e29b-41d4-a716-446655440000",
		Status: string(rules.UserStatusActive),
	}
}

func DeleteUserRequestDTO() string {
	return "550e8400-e29b-41d4-a716-446655440000"
}

func DeleteUserRequestInvalidDTO() string {
	return "invalid-uuid"
}

func GetUserRequestDTO() string {
	return "550e8400-e29b-41d4-a716-446655440000"
}

func GetUserRequestInvalidDTO() string {
	return "invalid-uuid"
}

func GetUserByEmailRequestDTO() string {
	return "test@example.com"
}

func GetUserByEmailNotFoundDTO() string {
	return "nonexistent@example.com"
}
