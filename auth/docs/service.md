# Service documentation

Auth service has 3 main services

- Auth
- Verification
- Email

### Auth

Auth service is used for registration, login, jwt tokens refreshing

```go
type AuthService interface {
    // register user with given email, password, role, returns jwt tokens
    RegisterUser(ctx context.Context, email string, password string, role string) (string, string, error)

    // login user with his email and password, returns jwt tokens
    LoginUser(ctx context.Context, email string, password string) (string, string, error)

    // refresh jwt tokens by given refresh token, returns jwt tokens
    RefreshTokens(ctx context.Context, refresh string) (string, string, error)

    // parse token and returns it's payload
    GetTokenPayload(ctx context.Context, access string) (*models.JwtPayload, error)

    // stop servicing
    Close()
}
```

### Verification

Verification service is used to verify user's account

```go
type VerificationService interface {
    // send verification email in email worker queue
    SendVerificationEmail(ctx context.Context, email string) error

    // verifies user's account by given token
    VerifyEmailByToken(ctx context.Context, token string) (string, string, error)

    // stop servicing
    Close()
}
```

### Email

Email service is used to send verification emails

```go
type EmailService interface {
    // send verification email with token
    SendVerificationEmail(email string, token string) error
}
```

### Email Producer

Email producer is a worker, that sends `message = {email, token}` in queue

```go
type EmailProducer interface {
    // send message in queue
    SendEmailInQueue(ctx context.Context, message *models.EmailMessage) error

    // stop servicing
    Close()
}
```

### Email Consumer

Email consumer is a worker, that reads messages from queue

```go
type EmailConsumer interface {
    // read message from queue and process it by fn
    ProcessEmailsFromQueue(ctx context.Context, fn func(msg *models.EmailMessage) error)

    // stop servicing
    Close()
}
```
