package user

import (
	"time"

	"github.com/IvanDrf/workk-hunter/pkg/users/internal/interfaces/grpc/dto"
	"github.com/google/uuid"
)

// user status
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBlocked  UserStatus = "blocked"
	UserStatusDeleted  UserStatus = "deleted"
)

// user role
type UserRole string

const (
	UserRoleUser      UserRole = "user"
	UserRoleModerator UserRole = "moderator"
	UserStatusAdmin   UserRole = "admin"
)

type User struct {
	ID          uuid.UUID `db:"id"`
	Username    string    `db:"username"`
	Email       string    `db:"email"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	PhoneNumber string    `db:"phone_number"`
	AvatarURL   string    `db:"avatar_url"`

	Status UserStatus `db:"status"`
	Role   UserRole   `db:"role"`

	Metadata map[string]string `db:"metadata"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updatet_at"`
}

func NewUser(req *dto.CreateUserRequest) *User {
	now := time.Now()
	return &User{
		ID:          uuid.New(),
		Username:    req.Username,
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,

		Status: UserStatusActive,
		Role:   UserRoleUser,

		Metadata: make(map[string]string),

		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) UpdateUser(req *dto.UpdateUserRequest) {
	if req.FirstName != "" {
		u.FirstName = req.FirstName
	}
	if req.LastName != "" {
		u.LastName = req.LastName
	}
	if req.PhoneNumber != "" {
		u.PhoneNumber = req.PhoneNumber
	}
	if req.AvatarURL != "" {
		u.AvatarURL = req.AvatarURL
	}

	if req.Metadata != nil {
		u.Metadata = req.Metadata
	}

	u.UpdatedAt = time.Now()
}
