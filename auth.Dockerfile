# Stage 2: Build authentication-service
FROM golang:1.23-alpine AS auth-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/task-ms

# Stage 3: Final image
FROM alpine:latest
WORKDIR /app
COPY --from=auth-builder /app/auth-service ./

EXPOSE 8084

# Install necessary tools
RUN apk update && apk add --no-cache mysql-client

# Entrypoint script (to handle secrets and start services)
COPY entrypoint.sh /app/
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh", "auth-service"]