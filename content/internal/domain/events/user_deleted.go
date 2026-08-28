package events

type UserDeletedEvent struct {
	UserID    string `json:"user_id"`
	DeletedAt string `jsin:"deleted_at"`
}
