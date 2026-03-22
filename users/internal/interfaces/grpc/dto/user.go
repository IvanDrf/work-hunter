package dto

// DTO for creating user
type CreateUserRequest struct {
	Username    string
	Email       string
	FirstName   string
	LastName    string
	PhoneNumber string
}

// DTO for updating user
type UpdateUserRequest struct {
	FirstName   string
	LastName    string
	PhoneNumber string
	AvatarURL   string
	Metadata    map[string]string
}
