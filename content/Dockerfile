FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o content_service ./cmd/api/main.go

FROM alpine:3.18

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/content_service .
COPY --from=builder /app/entry.sh .

RUN chmod +x entry.sh content_service

EXPOSE 8080

ENTRYPOINT ["./entry.sh"]