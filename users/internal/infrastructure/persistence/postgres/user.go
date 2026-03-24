package postgres

import (
	"context"
	"fmt"

	"github.com/IvanDrf/workk-hunter/pkg/users/internal/domain/models"
	"github.com/IvanDrf/workk-hunter/pkg/users/internal/interfaces/grpc/dto"
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

func isUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}
