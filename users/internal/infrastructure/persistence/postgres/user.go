package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
		id, username, email,
		first_name, last_name, phone_number,
		status, role, metadata, created_at, updated_at
	) VALUES(
		:id, :username, :email,
		:first_name, :last_name, :phone_number,
		:status, :role, :metadata, :created_at, :updated_at
	)`

	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		if isUniqueViolation(err) {
			return &models.Error{
				Message: fmt.Sprintf("user %v already exists", user),
				Code:    models.ErrCodeUserAlreadyExists,
			}
		}
		return &models.Error{
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
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &models.Error{
				Message: fmt.Sprintf("user with id = %v not found", id),
				Code:    models.ErrCodeUserNotFound,
			}
		}

		return nil, &models.Error{
			Message: fmt.Sprintf("failed to get user by id: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
	SELECT * FROM users
	WHERE username = $1 AND status != 'deleted'
	`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &models.Error{
				Message: fmt.Sprintf("user with username = %v not found", username),
				Code:    models.ErrCodeUserNotFound,
			}
		}

		return nil, &models.Error{
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
        phone_number = :phone_number,
        avatar_url = :avatar_url,
		metadata = :metadata,
        updated_at = :updated_at
	WHERE id = :id AND status != 'deleted'
	`

	_, err := r.db.NamedExecContext(ctx, query, updatedUser)
	if err != nil {
		return &models.Error{
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

	_, err := r.db.ExecContext(ctx, query, string(rules.UserStatusDeleted), time.Now(), id)

	if err != nil {
		return &models.Error{
			Message: fmt.Sprintf("failed to delete user: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}

func (r *UserRepository) ListUsers(ctx context.Context, params map[string]string) ([]*models.User, int32, error) {
	baseQuery := `
	SELECT * FROM users
	WHERE status != 'deleted'
	`

	countQuery := `SELECT COUNT(*) FROM users WHERE status != 'deleted'`

	var whereConditions string

	var offset string
	var limit string
	var orderBy string

	args := make([]string, 0, len(params))
	argPos := 1

	for key, val := range params {
		switch key {
		case "offset":
			offset = val

		case "limit":
			limit = val

		case "order_by":
			orderBy = val

		default:
			whereConditions += fmt.Sprintf(" AND %s = '$%d'", key, argPos)
			argPos++
			args = append(args, val)
		}
	}

	var totalCount int32
	if err := r.db.GetContext(ctx, &totalCount, countQuery+whereConditions, args); err != nil {
		return nil, 0, &models.Error{
			Message: fmt.Sprintf("failed to count users: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	if orderBy == "" {
		orderBy = "created_at"
	}
	if limit == "" || limit == "0" {
		limit = "100"
	}
	if offset == "" {
		offset = "0"
	}

	baseQuery += whereConditions
	baseQuery += fmt.Sprintf(" ORDER BY $%d LIMIT $%d OFFSET $%d", argPos, argPos+1, argPos+2)
	args = append(args, orderBy, limit, offset)

	var users []*models.User
	if err := r.db.SelectContext(ctx, &users, baseQuery, args); err != nil {
		return nil, 0, &models.Error{
			Message: fmt.Sprintf("failed to list users: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	return users, totalCount, nil
}

func (r *UserRepository) UpdateUserStatus(ctx context.Context, id uuid.UUID, status rules.UserStatus) error {
	query := `
	UPDATE users SET
		status = $1,
		updated_at = $2
	WHERE id = $3 AND status != 'deleted'
	`

	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return &models.Error{
			Message: fmt.Sprintf("failed to update user status: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}

func isUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}
