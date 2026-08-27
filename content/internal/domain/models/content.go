package models

type ContentType string

const (
	TypeResume ContentType = "resume"
	TypeAvatar ContentType = "avatar"
)

type ContentMetadata struct {
	UserID   string
	Type     ContentType
	MimeType string //MIME (application/pdf, image/jpeg)
	Size     int64
}
