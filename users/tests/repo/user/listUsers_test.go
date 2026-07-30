package user

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/internal/domain/rules"
	"github.com/IvanDrf/work-hunter/users/tests/repo/user/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestListUsers_Success(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	users := fixtures.CreateUsers()
	params := fixtures.ListUsersParams()

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(len(users))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
		WillReturnRows(countRows)

	rows := sqlmock.NewRows([]string{
		"id", "email", "first_name", "last_name", "company_name",
		"status", "role", "verificated", "created_at", "updated_at",
	})
	for _, user := range users {
		rows.AddRow(
			user.ID, user.Email, user.FirstName, user.LastName, user.CompanyName,
			user.Status, user.Role, user.Verificated, user.CreatedAt, user.UpdatedAt,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM users").
		WillReturnRows(rows)

	result, totalCount, err := repo.ListUsers(t.Context(), params)

	assert.NoError(t, err)
	assert.Len(t, result, len(users))
	assert.Equal(t, int32(len(users)), totalCount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListUsers_WithStatusFilter(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	users := fixtures.CreateUsers()
	params := &models.ListUsersParams{
		Page:     1,
		PageSize: 10,
		Status:   string(rules.UserStatusActive),
	}

	// Фильтруем пользователей по статусу
	filteredUsers := make([]*models.User, 0)
	for _, u := range users {
		if u.Status == rules.UserStatusActive {
			filteredUsers = append(filteredUsers, u)
		}
	}

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(len(filteredUsers))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
		WillReturnRows(countRows)

	rows := sqlmock.NewRows([]string{
		"id", "email", "first_name", "last_name", "company_name",
		"status", "role", "verificated", "created_at", "updated_at",
	})
	for _, user := range filteredUsers {
		rows.AddRow(
			user.ID, user.Email, user.FirstName, user.LastName, user.CompanyName,
			user.Status, user.Role, user.Verificated, user.CreatedAt, user.UpdatedAt,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM users").
		WillReturnRows(rows)

	result, totalCount, err := repo.ListUsers(t.Context(), params)

	assert.NoError(t, err)
	assert.Len(t, result, len(filteredUsers))
	assert.Equal(t, int32(len(filteredUsers)), totalCount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListUsers_WithRoleFilter(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	users := fixtures.CreateUsers()
	params := &models.ListUsersParams{
		Page:     1,
		PageSize: 10,
		Role:     string(rules.UserRoleAdmin),
	}

	filteredUsers := make([]*models.User, 0)
	for _, u := range users {
		if u.Role == rules.UserRoleAdmin {
			filteredUsers = append(filteredUsers, u)
		}
	}

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(len(filteredUsers))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
		WillReturnRows(countRows)

	rows := sqlmock.NewRows([]string{
		"id", "email", "first_name", "last_name", "company_name",
		"status", "role", "verificated", "created_at", "updated_at",
	})
	for _, user := range filteredUsers {
		rows.AddRow(
			user.ID, user.Email, user.FirstName, user.LastName, user.CompanyName,
			user.Status, user.Role, user.Verificated, user.CreatedAt, user.UpdatedAt,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM users").
		WillReturnRows(rows)

	result, totalCount, err := repo.ListUsers(t.Context(), params)

	assert.NoError(t, err)
	assert.Len(t, result, len(filteredUsers))
	assert.Equal(t, int32(len(filteredUsers)), totalCount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListUsers_WithSearchQuery(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	user := fixtures.CreateSingleUser()
	params := &models.ListUsersParams{
		Page:        1,
		PageSize:    10,
		SearchQuery: "john",
	}

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
		WillReturnRows(countRows)

	rows := sqlmock.NewRows([]string{
		"id", "email", "first_name", "last_name", "company_name",
		"status", "role", "verificated", "created_at", "updated_at",
	}).AddRow(
		user.ID, user.Email, user.FirstName, user.LastName, user.CompanyName,
		user.Status, user.Role, user.Verificated, user.CreatedAt, user.UpdatedAt,
	)

	mock.ExpectQuery("SELECT \\* FROM users").
		WillReturnRows(rows)

	result, totalCount, err := repo.ListUsers(t.Context(), params)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int32(1), totalCount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListUsers_WithCombinedFilters(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	users := fixtures.CreateUsers()
	params := &models.ListUsersParams{
		Page:        1,
		PageSize:    10,
		Status:      string(rules.UserStatusActive),
		Role:        string(rules.UserRoleEmployee),
		SearchQuery: "john",
	}

	filteredUsers := make([]*models.User, 0)
	for _, u := range users {
		if u.Status == rules.UserStatusActive &&
			u.Role == rules.UserRoleEmployee &&
			(u.FirstName == "John" || u.LastName == "Doe" || u.Email == "john.doe@example.com") {
			filteredUsers = append(filteredUsers, u)
		}
	}

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(len(filteredUsers))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
		WillReturnRows(countRows)

	rows := sqlmock.NewRows([]string{
		"id", "email", "first_name", "last_name", "company_name",
		"status", "role", "verificated", "created_at", "updated_at",
	})
	for _, user := range filteredUsers {
		rows.AddRow(
			user.ID, user.Email, user.FirstName, user.LastName, user.CompanyName,
			user.Status, user.Role, user.Verificated, user.CreatedAt, user.UpdatedAt,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM users").
		WillReturnRows(rows)

	result, totalCount, err := repo.ListUsers(t.Context(), params)

	assert.NoError(t, err)
	assert.Len(t, result, len(filteredUsers))
	assert.Equal(t, int32(len(filteredUsers)), totalCount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListUsers_EmptyResult(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	params := fixtures.ListUsersParams()

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
		WillReturnRows(countRows)

	result, totalCount, err := repo.ListUsers(t.Context(), params)

	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int32(0), totalCount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListUsers_CountError(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	params := fixtures.ListUsersParams()

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
		WillReturnError(errors.New("count query failed"))

	result, totalCount, err := repo.ListUsers(t.Context(), params)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, int32(0), totalCount)
	assert.Contains(t, err.Error(), "failed to count users")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListUsers_SelectError(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	params := fixtures.ListUsersParams()

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
		WillReturnRows(countRows)

	mock.ExpectQuery("SELECT \\* FROM users").
		WillReturnError(errors.New("select query failed"))

	result, totalCount, err := repo.ListUsers(t.Context(), params)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, int32(0), totalCount)
	assert.Contains(t, err.Error(), "failed to list users")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListUsers_WithPagination(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	users := fixtures.CreateUsers()
	params := &models.ListUsersParams{
		Page:     2,
		PageSize: 2,
	}

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(len(users))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
		WillReturnRows(countRows)

	rows := sqlmock.NewRows([]string{
		"id", "email", "first_name", "last_name", "company_name",
		"status", "role", "verificated", "created_at", "updated_at",
	})
	startIdx := (params.Page - 1) * params.PageSize
	for i := startIdx; i < len(users) && i < startIdx+params.PageSize; i++ {
		user := users[i]
		rows.AddRow(
			user.ID, user.Email, user.FirstName, user.LastName, user.CompanyName,
			user.Status, user.Role, user.Verificated, user.CreatedAt, user.UpdatedAt,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM users").
		WillReturnRows(rows)

	result, totalCount, err := repo.ListUsers(t.Context(), params)

	assert.NoError(t, err)
	expectedCount := 1
	assert.Len(t, result, expectedCount)
	assert.Equal(t, int32(3), totalCount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListUsers_PageSizeMaxLimit(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	params := &models.ListUsersParams{
		Page:     1,
		PageSize: 150,
	}
	users := fixtures.CreateUsers()

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(len(users))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
		WillReturnRows(countRows)

	rows := sqlmock.NewRows([]string{
		"id", "email", "first_name", "last_name", "company_name",
		"status", "role", "verificated", "created_at", "updated_at",
	})
	for _, user := range users {
		rows.AddRow(
			user.ID, user.Email, user.FirstName, user.LastName, user.CompanyName,
			user.Status, user.Role, user.Verificated, user.CreatedAt, user.UpdatedAt,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM users").
		WillReturnRows(rows)

	result, totalCount, err := repo.ListUsers(t.Context(), params)

	assert.NoError(t, err)
	assert.Len(t, result, len(users))
	assert.Equal(t, int32(len(users)), totalCount)
	assert.Equal(t, 100, params.PageSize)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListUsers_WithSorting(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	users := fixtures.CreateUsers()
	params := &models.ListUsersParams{
		Page:     1,
		PageSize: 10,
		SortBy:   "email",
		SortDesc: true,
	}

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(len(users))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
		WillReturnRows(countRows)

	rows := sqlmock.NewRows([]string{
		"id", "email", "first_name", "last_name", "company_name",
		"status", "role", "verificated", "created_at", "updated_at",
	})
	for _, user := range users {
		rows.AddRow(
			user.ID, user.Email, user.FirstName, user.LastName, user.CompanyName,
			user.Status, user.Role, user.Verificated, user.CreatedAt, user.UpdatedAt,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM users").
		WillReturnRows(rows)

	result, totalCount, err := repo.ListUsers(t.Context(), params)

	assert.NoError(t, err)
	assert.Len(t, result, len(users))
	assert.Equal(t, int32(len(users)), totalCount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListUsers_InvalidPageDefaults(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	users := fixtures.CreateUsers()
	params := &models.ListUsersParams{
		Page:     0, // Invalid
		PageSize: 0, // Invalid
	}

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(len(users))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
		WillReturnRows(countRows)

	rows := sqlmock.NewRows([]string{
		"id", "email", "first_name", "last_name", "company_name",
		"status", "role", "verificated", "created_at", "updated_at",
	})
	for _, user := range users {
		rows.AddRow(
			user.ID, user.Email, user.FirstName, user.LastName, user.CompanyName,
			user.Status, user.Role, user.Verificated, user.CreatedAt, user.UpdatedAt,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM users").
		WillReturnRows(rows)

	result, totalCount, err := repo.ListUsers(t.Context(), params)

	assert.NoError(t, err)
	assert.Len(t, result, len(users))
	assert.Equal(t, int32(len(users)), totalCount)
	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 10, params.PageSize)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListUsers_ExcludeDeleted(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	users := fixtures.CreateUsers()

	params := fixtures.ListUsersParams()

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(len(users))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
		WillReturnRows(countRows)

	rows := sqlmock.NewRows([]string{
		"id", "email", "first_name", "last_name", "company_name",
		"status", "role", "verificated", "created_at", "updated_at",
	})
	for _, user := range users {
		rows.AddRow(
			user.ID, user.Email, user.FirstName, user.LastName, user.CompanyName,
			user.Status, user.Role, user.Verificated, user.CreatedAt, user.UpdatedAt,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM users").
		WillReturnRows(rows)

	result, totalCount, err := repo.ListUsers(t.Context(), params)

	assert.NoError(t, err)
	assert.Len(t, result, len(users))
	assert.Equal(t, int32(len(users)), totalCount)
	for _, u := range result {
		assert.NotEqual(t, rules.UserStatusDeleted, u.Status)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}
