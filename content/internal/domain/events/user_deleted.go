package events

type UserDeletedEvent struct {
	UserID string `json:"user_id"`
}
