# Stage 1: Build task-service
FROM golang:1.22-alpine AS task-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o task-service ./cmd/task-ms

# Stage 2: Build authentication-service
FROM golang:1.22-alpine AS auth-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/task-ms

# Stage 3: Final image
FROM alpine:latest
WORKDIR /app
COPY --from=task-builder /app/task-service ./
COPY --from=auth-builder /app/auth-service ./

# Install necessary tools
RUN apk update && apk add --no-cache mysql-client

# Expose ports
EXPOSE 3306
EXPOSE 8083
EXPOSE 8084

# Entrypoint script (to handle secrets and start services)
COPY entrypoint.sh /app/
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]