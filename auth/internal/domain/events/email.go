package events

type EmailWorker interface {
	SendEmailVerification(email string)
}
