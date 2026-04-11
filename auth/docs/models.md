# Models documentation

### Errors

The model is responsible for business-logic errors

- Message
- Code

```go
type Error struct {
    Message string    `json:"message"`
    Code    ErrorCode `json:"code"`
}
```

Codes

```go
type ErrorCode string

const (
    // internal errors, example: can't generate jwt tokens
    ErrCodeInternal ErrorCode = "INTERNAL_ERROR"

    // external error, not supported role was given in resquest
    ErrCodeInvalidUserRole        ErrorCode = "INVALID_USER_ROLE"
    // can't find user with given parameter
    ErrCodeUserNotFound           ErrorCode = "USER_NOT_FOUND"
    // user already exists, can't create with the same email
    ErrCodeUserAlreadyExists      ErrorCode = "USER_ALREADY_EXISTS"
    // user has already status "verificated", there is no need to send verification email again
    ErrCodeUserAlreadyVerificated ErrorCode = "USER_ALREADY_VERIFICATED"
    // verification token already outdated
    ErrCodeOutdatedToken          ErrorCode = "TOKEN_IS_OUTDATED"

    // incorrect password for user was in request
    ErrCodeInvalidPassword ErrorCode = "INVALID_PASSWORD"
    // invalid email format
    ErrCodeInvalidEmail    ErrorCode = "INVALID_EMAIL"

    // invalid jwt token
    ErrCodeInvalidJWT ErrorCode = "INVALID_JWT_TOKENS"
)
```

### User

The model consists user's information

- ID
- Email
- HashedPassword
- user's role
- verificated status

```go
type User struct {
    ID uuid.UUID `json:"user_id"`

    Email          string `json:"email"`
    HashedPassword string `json:"password"`
    Role           Role   `json:"role"`

    Verificated bool `json:"verificated"`
}
```

### Role

The model is responsible for user's role

```go
type Role string

const (
    ADMIN    Role = "ADMIN"
    EMPLOYEE Role = "EMPLOYEE" // can respond to vacancies
    EMPLOYER Role = "EMPLOYER" // can create vacancies
)
```

### Jwt

The model is responsible for jwt payload

- user's ID
- verificated status
- user's role

```go
type JwtPayload struct {
    UserID      string `json:"user_id"`
    Verificated bool   `json:"verificated"`
    Role        Role   `json:"role"`
}
```

### Message

The model is responsible for messages in email verification queue

- email
- verification token
- expiring time for verification token

```go
type EmailMessage struct {
    Email string `json:"email"`

    Token string    `json:"token"`
    Exp   time.Time `json:"exp"`
}
```
