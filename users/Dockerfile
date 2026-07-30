FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/migrator ./cmd/migrator/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /bin/server ./server
COPY --from=builder /bin/migrator ./migrator
COPY --from=builder /app/entry.sh ./entry.sh

COPY --from=builder /app/internal/infrastructure/persistence/postgres/migrations /app/migrations

RUN chmod +x ./entry.sh

EXPOSE 8080

ENTRYPOINT ["./entry.sh"]