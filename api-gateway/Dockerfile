FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api_gateway ./cmd/main.go

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/api_gateway .
COPY --from=builder /app/.env .

EXPOSE 8000
CMD ["./api_gateway"]