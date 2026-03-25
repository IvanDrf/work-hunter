package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/IvanDrf/workk-hunter/pkg/users/internal/domain/models"
	"github.com/IvanDrf/workk-hunter/pkg/users/internal/domain/rules"
	"github.com/IvanDrf/workk-hunter/pkg/users/internal/interfaces/grpc/dto"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(conn *PostgresConnection) *UserRepository {
	return &UserRepository{
		db: conn.GetDB(),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *dto.CreateUserRequest) (*models.User, error) {
	query := `
	INSERT INTO users (
		id, username, email,
		first_name, last_name, phone_number,
		status, role, metadata, created_at, updated_at
	) VALUES(
		:id, :username, :email,
		:first_name, :last_name, :phone_number,
		:status, :role, :metadata, :created_at, :updated_at
	)
	RETURNING *`

	userModel := models.NewUser(user)

	rows, err := r.db.NamedQueryContext(ctx, query, userModel)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, &models.Error{
				Message: fmt.Sprintf("user %v already exists", userModel),
				Code:    models.ErrCodeUserAlreadyExists,
			}
		}
		return nil, &models.Error{
			Message: fmt.Sprintf("failed to create user %v: %v", userModel, err),
			Code:    models.ErrCodeInternal,
		}
	}

	defer rows.Close()

	var created models.User
	if rows.Next() {
		if err := rows.StructScan(&created); err != nil {
			return nil, &models.Error{
				Message: fmt.Sprintf("failed to scan created user: %v", err),
				Code:    models.ErrCodeInternal,
			}
		}
	}

	return &created, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
	SELECT * FROM users
	WHERE id = $1 AND status != 'deleted
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
	WHERE username = $1 AND status != 'deleted
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

func (r *UserRepository) UpdateUser(ctx context.Context, req *dto.UpdateUserRequest) (*models.User, error) {
	current, err := r.GetUserByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	current.UpdateUser(req)

	query := `
	UPDATE users SET
		first_name = :first_name,
		last_name = :last_name,
        phone_number = :phone_number,
        avatar_url = :avatar_url,
		metadata = :metadata,
        updated_at = :updated_at
	WHERE id = :id AND status != 'deleted'
	RETURNING *
	`

	rows, err := r.db.NamedQueryContext(ctx, query, current)
	if err != nil {
		return nil, &models.Error{
			Message: fmt.Sprintf("failed to update user: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	defer rows.Close()

	var updated models.User
	if rows.Next() {
		if err := rows.StructScan(&updated); err != nil {
			return nil, &models.Error{
				Message: fmt.Sprintf("failed to scan updated user: %v", err),
				Code:    models.ErrCodeInternal,
			}
		}
	}

	return &updated, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	checkQuery := `
	SELECT status FROM users
	WHERE id = $1
	`
	var status string
	err := r.db.QueryRowContext(ctx, checkQuery, id).Scan(&status)
	if err != nil {
		return &models.Error{
			Message: "failed to check user",
			Code:    models.ErrCodeInternal,
		}
	}

	if status == "deleted" {
		query := `
		DELETE FROM users WHERE id = $1
		`
		_, err = r.db.ExecContext(ctx, query, id)
	} else {
		query := `
		UPDATE users SET
			status = $1,
			updated_at = $2,
			deleted_at = $2
		WHERE id = $3 AND status != 'deleted'
		`

		_, err = r.db.ExecContext(ctx, query, string(rules.UserStatusDeleted), time.Now(), id)
	}

	if err != nil {
		return &models.Error{
			Message: fmt.Sprintf("failed to delete user: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}

func (r *UserRepository) ListUsers(ctx context.Context, req *dto.ListUsersRequest) ([]*models.User, int32, error) {
	//TODO
	return nil, 0, nil
}

func (r *UserRepository) UpdateUserStatus(ctx context.Context, req *dto.UpdateUserStatusRequest) (*models.User, error) {
	if err := rules.ValidateUserStatus(req.Status); err != nil {
		return nil, &models.Error{
			Message: "invalid user status: " + string(req.Status),
			Code:    models.ErrCodeInvalidStatus,
		}
	}

	query := `
	UPDATE users SET
		status = $1,
		updated_at = $2
	WHERE id = $3 AND status != 'deleted'
	RETURNING *
	`

	var updated models.User
	err := r.db.GetContext(ctx, &updated, query, req.Status, time.Now(), req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &models.Error{
				Message: "user not found",
				Code:    models.ErrCodeUserNotFound,
			}
		}

		return nil, &models.Error{
			Message: fmt.Sprintf("failed to update user status: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	return &updated, nil
}

func isUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}
