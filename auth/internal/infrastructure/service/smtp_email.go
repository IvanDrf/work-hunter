package service

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/smtp"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
)

type SmtpEmailService struct {
	host string

	username string
	password string
}

func NewSmtpEmailService(host string, username string, password string) *SmtpEmailService {
	return &SmtpEmailService{
		host:     host,
		username: username,
		password: password,
	}
}

func (e *SmtpEmailService) SendVerificationEmail(email string, token string) error {
	message, err := createVerificationMessage(token)
	if err != nil {
		slog.Error("SendVerifEmail error", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't create verification message",
			Code:    models.ErrCodeInternal,
		}
	}

	auth := smtp.PlainAuth("", e.username, e.password, "smtp.gmail.com")
	err = smtp.SendMail(e.host, auth, e.username, []string{email}, []byte(message))
	if err != nil {
		slog.Error("SendVerifEmail error", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't send verification email",
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}

const (
	htmlEmailBodyPath = "static/email.html"
	headers           = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\nSubject: Подтверждение email на Work-Hunter\n"
)

func createVerificationMessage(token string) (string, error) {
	tmpl, err := template.ParseFiles(htmlEmailBodyPath)
	if err != nil {
		return "", err
	}

	buff := bytes.Buffer{}

	err = tmpl.Execute(&buff, struct {
		Token string
	}{Token: token},
	)

	if err != nil {
		return "", err
	}

	return headers + "\n\n" + buff.String(), nil
}
