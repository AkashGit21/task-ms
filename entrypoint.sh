#!/bin/sh

GO_ENV="local"

# Server related envs
TASK_APP_HOST=""
TASK_APP_PORT=8083

AUTH_APP_HOST=""
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

# Start services
APP_HOST="$TASK_APP_HOST" APP_PORT="$TASK_APP_PORT" TASK_DB_HOST="$TASK_DB_HOST" TASK_DB_PORT="$TASK_DB_PORT" TASK_DB_USER="$TASK_DB_USER" TASK_DB_PASSWORD="$TASK_DB_PASSWORD" TASK_DB_NAME="$TASK_DB_NAME" ./task-service task &
APP_HOST="$AUTH_APP_HOST" APP_PORT="$AUTH_APP_PORT" AUTH_DB_HOST="$AUTH_DB_HOST" AUTH_DB_PORT="$AUTH_DB_PORT" AUTH_DB_USER="$AUTH_DB_USER" AUTH_DB_PASSWORD="$AUTH_DB_PASSWORD" AUTH_DB_NAME="$AUTH_DB_NAME" ./auth-service authn


# Test MySQL connection with openssl
openssl s_client -connect mysql:3306 -showcerts 2>/dev/null | openssl x509 -text -noout

# Keep the container running
wait
