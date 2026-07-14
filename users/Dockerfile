FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/migrator ./cmd/migrator/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /bin/server /app/server
COPY --from=builder /bin/migrator /app/migrator

COPY --from=builder /app/internal/infrastructure/persistence/postgres/migrations /app/internal/infrastructure/persistence/postgres/migrations

RUN mkdir -p /app/config

CMD ["/app/server"]