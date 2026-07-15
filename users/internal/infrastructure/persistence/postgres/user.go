package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/internal/domain/rules"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserRepository struct {
	*PostgresConnection
}

func NewUserRepository(conn *PostgresConnection) *UserRepository {
	return &UserRepository{
		PostgresConnection: conn,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
	INSERT INTO users (
		id, email,
		first_name, last_name, company_name,
		status, role, verificated
	) VALUES(
		:id, :email,
		:first_name, :last_name, :company_name,
		:status, :role, :verificated
	)`

	_, err := r.DB.NamedExecContext(ctx, query, user)
	if err != nil {
		if isUniqueViolation(err) {
			return models.Error{
				Message: fmt.Sprintf("user %v already exists", user),
				Code:    models.ErrCodeUserAlreadyExists,
			}
		}
		return models.Error{
			Message: fmt.Sprintf("failed to create user %v: %v", user, err),
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
	SELECT * FROM users
	WHERE id = $1
	`

	var user models.User
	err := r.DB.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.Error{
				Message: fmt.Sprintf("user with id = %v not found", id),
				Code:    models.ErrCodeUserNotFound,
			}
		}

		return nil, models.Error{
			Message: fmt.Sprintf("failed to get user by id: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
	SELECT * FROM users
	WHERE email = $1 AND status != 'deleted'
	`

	var user models.User
	err := r.DB.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.Error{
				Message: fmt.Sprintf("user with email = %v not found", email),
				Code:    models.ErrCodeUserNotFound,
			}
		}

		return nil, models.Error{
			Message: fmt.Sprintf("failed to get user by username: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, updatedUser *models.User) error {
	query := `
	UPDATE users SET
		first_name = :first_name,
		last_name = :last_name,
        company_name = :company_name,
		verificated = :verificated
	WHERE id = :id AND status != 'deleted'
	`

	_, err := r.DB.NamedExecContext(ctx, query, updatedUser)
	if err != nil {
		return models.Error{
			Message: fmt.Sprintf("failed to update user: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID, permanent bool) error {
	var query string
	if permanent {
		query = `
		DELETE FROM users
		WHERE id = $1;
		`
	} else {
		query = `
		UPDATE users SET status = 'deleted' WHERE id = $1 AND status != 'deleted'
		`
	}

	_, err := r.DB.ExecContext(ctx, query, id)

	if err != nil {
		return models.Error{
			Message: fmt.Sprintf("failed to delete user: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}

func (r *UserRepository) ListUsers(ctx context.Context, params *models.ListUsersParams) ([]*models.User, int32, error) {
	baseQuery := `FROM users WHERE status != 'deleted'`

	var conditions []string
	var args []any
	argPos := 1

	if params.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argPos))
		args = append(args, params.Status)
		argPos++
	}

	if params.Role != "" {
		conditions = append(conditions, fmt.Sprintf("role = $%d", argPos))
		args = append(args, params.Role)
		argPos++
	}

	if params.SearchQuery != "" {
		searchPattern := "%" + params.SearchQuery + "%"
		conditions = append(conditions, fmt.Sprintf(
			"(email ILIKE $%d OR first_name ILIKE $%d OR last_name ILIKE $%d)",
			argPos, argPos+1, argPos+2,
		))
		args = append(args, searchPattern, searchPattern, searchPattern)
		argPos += 3
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " AND " + strings.Join(conditions, " AND ")
	}

	orderClause := " ORDER BY created_at DESC"
	if params.SortBy != "" {
		direction := "ASC"
		if params.SortDesc {
			direction = "DESC"
		}
		orderClause = fmt.Sprintf(" ORDER BY %s %s", params.SortBy, direction)
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}

	offset := (params.Page - 1) * params.PageSize
	paginationClause := fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, params.PageSize, offset)

	countQuery := "SELECT COUNT(*) " + baseQuery + whereClause
	var totalCount int32
	if err := r.DB.GetContext(ctx, &totalCount, countQuery, args[:len(args)-2]...); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	if totalCount == 0 {
		return []*models.User{}, 0, nil
	}

	query := "SELECT * " + baseQuery + whereClause + orderClause + paginationClause
	var users []*models.User
	if err := r.DB.SelectContext(ctx, &users, query, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, totalCount, nil
}

func (r *UserRepository) UpdateUserStatus(ctx context.Context, id uuid.UUID, status rules.UserStatus) error {
	query := `
	UPDATE users SET
		status = $1
	WHERE id = $2;
	`

	_, err := r.DB.ExecContext(ctx, query, status, id)
	if err != nil {
		return models.Error{
			Message: fmt.Sprintf("failed to update user status: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}

func (r *UserRepository) Close() {
	r.PostgresConnection.Close()
}

func isUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}
