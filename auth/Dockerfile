FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o auth_service ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o migrator ./cmd/migrator/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o email_worker ./cmd/worker/main.go

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/auth_service .
COPY --from=builder /app/migrator .
COPY --from=builder /app/email_worker .
COPY --from=builder /app/.env .
COPY --from=builder /app/entry.sh .
RUN chmod +x entry.sh 
COPY --from=builder /app/internal/infrastructure/persistence/postgres/migrations/ .

EXPOSE 50051
ENTRYPOINT ["./entry.sh"]