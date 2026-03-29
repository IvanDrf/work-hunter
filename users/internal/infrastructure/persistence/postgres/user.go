package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/internal/domain/rules"
	"github.com/IvanDrf/work-hunter/users/internal/interfaces/grpc/dto"
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

	userModel := models.NewUser(user.ID, user.Username, user.Email, user.FirstName, user.LastName, user.PhoneNumber)

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

	current.UpdateUser(req.FirstName, req.LastName, req.PhoneNumber, req.AvatarURL, req.Metadata)

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

func (r *UserRepository) ListUsers(ctx context.Context, req *dto.ListUsersRequest) (*dto.ListUsersResponse, error) {
	baseQuery := `
	SELECT * FROM users 
	WHERE status != 'deleted'
	`

	countQuery := "SELECT COUNT(*) FROM users WHERE status != 'deleted'"

	conditions := []string{}
	args := []any{}
	argPos := 1

	if req.Status != "" {
		if err := rules.ValidateUserStatus(req.Status); err != nil {
			return nil, &models.Error{
				Message: "invalid user status in request",
				Code:    models.ErrCodeInvalidStatus,
			}
		}
		conditions = append(conditions, fmt.Sprintf("status = $%d", argPos))
		args = append(args, string(req.Status))
		argPos++
	}

	if req.Role != "" {
		if err := rules.ValidateUserRole(req.Role); err != nil {
			return nil, &models.Error{
				Message: "invalod user role in request",
				Code:    models.ErrCodeInvalidRole,
			}
		}
		conditions = append(conditions, fmt.Sprintf("role = $%d", argPos))
		args = append(args, string(req.Role))
		argPos++
	}

	if req.SearchQuery != "" {
		searchPattern := "%" + req.SearchQuery + "%"
		conditions = append(conditions, fmt.Sprintf(
			"(username ILIKE $%d OR email ILIKE $%d OR first_name ILIKE $%d OR last_name ILIKE $%d)",
			argPos, argPos, argPos, argPos,
		))
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern)
		argPos += 4
	}

	if len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	var totalCount int32
	err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, &models.Error{
			Message: fmt.Sprintf("failed to get total count: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	sortField := "created_at"
	if req.SortBy != "" {
		allowedSortFields := map[string]bool{
			"id": true, "username": true, "email": true,
			"created_at": true, "updated_at": true, "last_login_at": true,
		}
		if allowedSortFields[req.SortBy] {
			sortField = req.SortBy
		}
	}

	query := fmt.Sprintf("%s ORDER BY %s ASC LIMIT $%d OFFSET $%d",
		baseQuery, sortField, argPos, argPos+1)
	args = append(args, pageSize, req.Offset)

	var users []*models.User
	if err = r.db.SelectContext(ctx, &users, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, &models.Error{
				Message: "no users with such parametrs",
				Code:    models.ErrCodeUserNotFound,
			}
		}

		return nil, &models.Error{
			Message: fmt.Sprintf("failed to get users: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	hasNext := int32(len(users)) == pageSize && (req.Offset+pageSize) < totalCount

	return &dto.ListUsersResponse{
		Users:      users,
		TotalCount: totalCount,
		HasNext:    hasNext,
	}, nil
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
