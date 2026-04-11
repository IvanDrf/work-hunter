# Repo documentation

Service has 2 main repos:

- User Repo
- Token Repo

### User Repo

User repo working with users data: email, password _(in database password is hashed)_, role, verificated status

Methods:

```go
type UserRepo interface {
    // create user in database
    CreateUser(ctx context.Context, user *models.User) error
    // find user in database by his email
    FindUserByEmail(ctx context.Context, email string) (*models.User, error)

    // set verificated status = true by user's email
    VerifyEmail(ctx context.Context, email string) error

    // closing connection to database
    Close()
}
```

### Token Repo

Token repo working with verification tokens, that are used in email verification

Methods:

```go
type TokenRepo interface {
    // create new token in database
    CreateToken(ctx context.Context, email string, token string, ttl time.Duration) error

    // find user's email by verification token
    FindEmailByToken(ctx context.Context, token string) string

    // delete token in databse
    DeleteToken(ctx context.Context, token string) error

    // close connection to database
    Close()
}
```
