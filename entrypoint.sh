#!/bin/sh

SERVICE="${1:-all}"

GO_ENV="local"

# Server related envs
TASK_APP_HOST=${TASK_APP_HOST:-"0.0.0.0"}
TASK_APP_PORT=8083

AUTH_APP_HOST=${AUTH_APP_HOST:-"0.0.0.0"}
AUTH_APP_PORT=8084

# Database envs
TASK_DB_HOST=${TASK_DB_HOST:-"mysql"}
TASK_DB_PORT=${TASK_DB_PORT:-"3306"}
TASK_DB_USER=${TASK_DB_USER:-"task_user"}
TASK_DB_PASSWORD=${TASK_DB_PASSWORD:-"task_password"}
TASK_DB_NAME=${TASK_DB_NAME:-"task_db"}

AUTH_DB_HOST=${AUTH_DB_HOST:-"mysql"}
AUTH_DB_PORT=${AUTH_DB_PORT:-"3306"}
AUTH_DB_USER=${AUTH_DB_USER:-"auth_user"}
AUTH_DB_PASSWORD=${AUTH_DB_PASSWORD:-"auth_password"}
AUTH_DB_NAME=${AUTH_DB_NAME:-"auth_db"}

echo "Service is $SERVICE"
# Start services conditionally
if [ "$SERVICE" = "task-service" ] || [ "$SERVICE" = "all" ]; then
    APP_LOG_LEVEL="INFO" APP_HOST="$TASK_APP_HOST" APP_PORT="$TASK_APP_PORT" TASK_DB_HOST="$TASK_DB_HOST" TASK_DB_PORT="$TASK_DB_PORT" TASK_DB_USER="$TASK_DB_USER" TASK_DB_PASSWORD="$TASK_DB_PASSWORD" TASK_DB_NAME="$TASK_DB_NAME" AUTH_JWT_SECRET="random secret value" ./task-service task
fi

if [ "$SERVICE" = "auth-service" ] || [ "$SERVICE" = "all" ]; then
    APP_LOG_LEVEL="INFO" APP_HOST="$AUTH_APP_HOST" APP_PORT="$AUTH_APP_PORT" AUTH_DB_HOST="$AUTH_DB_HOST" AUTH_DB_PORT="$AUTH_DB_PORT" AUTH_DB_USER="$AUTH_DB_USER" AUTH_DB_PASSWORD="$AUTH_DB_PASSWORD" AUTH_DB_NAME="$AUTH_DB_NAME" AUTH_JWT_SECRET="random secret value" ./auth-service authn
fi

# Test MySQL connection with openssl
openssl s_client -connect mysql:3306 -showcerts 2>/dev/null | openssl x509 -text -noout

# Keep the container running
wait
