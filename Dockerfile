# Build stage
FROM golang:1.26.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bot .

# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bot .

CMD ["./bot"]