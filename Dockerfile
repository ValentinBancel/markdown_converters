FROM golang:1.24.7-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o fiber-app main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/fiber-app .
COPY --from=builder /app/.env.example .env

EXPOSE 3000
CMD ["./fiber-app"]