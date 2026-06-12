# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/server/main.go

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
# Install CA certificates for Firebase/HTTPS
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/configs/firebase-service-account.json ./configs/
EXPOSE 8080
CMD ["./main"]
