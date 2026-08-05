package rules

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/IvanDrf/work-hunter/content/internal/domain/models"
)

const (
	MaxResumeSize = 5 * 1024 * 1024 // 5 MB
	MaxAvatarSize = 2 * 1024 * 1024 // 2 MB
)

// AllowedMimeTypes defines the allowed MIME types for each content type
var AllowedMimeTypes = map[models.ContentType]map[string]struct{}{
	models.TypeResume: {
		"application/pdf": struct{}{},
	},
	models.TypeAvatar: {
		"image/png":  struct{}{},
		"image/jpeg": struct{}{},
	},
}

// ValidatedFile contains the restored Reader and the verified MIME type
type ValidatedFile struct {
	Reader   io.Reader
	MimeType string
	Size     int64
}

// ValidateContent checks the size and the true MIME type of the file based on its content (Magic Bytes).
// Since reading the header shifts the io.Reader pointer, the function returns io.MultiReader.
// for the subsequent correct sending of the file to S3.
func ValidateContent(cType models.ContentType, file io.Reader, size int64) (*ValidatedFile, error) {
	if size <= 0 {
		return nil, models.ErrFileEmpty
	}

	if err := validateSize(cType, size); err != nil {
		return nil, err
	}

	head := make([]byte, 512)
	n, err := file.Read(head)
	if err != nil && err != io.EOF {
		return nil, models.ErrReadFailed
	}

	detectedMime := http.DetectContentType(head[:n])

	allowedTypes, typeExists := AllowedMimeTypes[cType]
	_, allowed := allowedTypes[detectedMime]
	if !typeExists || !allowed {
		return nil, models.ErrInvalidMimeType
	}

	restoredFile := io.MultiReader(bytes.NewReader(head[:n]), file)

	return &ValidatedFile{
		Reader:   restoredFile,
		MimeType: detectedMime,
		Size:     size,
	}, nil
}

func validateSize(cType models.ContentType, size int64) error {
	switch cType {
	case models.TypeResume:
		if size > MaxResumeSize {
			return models.ErrFileTooLarge
		}

	case models.TypeAvatar:
		if size > MaxAvatarSize {
			return models.ErrFileTooLarge
		}

	default:
		return errors.New("unknown content type")
	}

	return nil
}
