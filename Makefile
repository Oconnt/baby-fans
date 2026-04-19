.PHONY: build docker-compose-up docker-compose-down docker-compose-restart docker-shell docker-push migrate nginx-up nginx-down nginx-logs nginx-restart

# Build variables
BINARY_NAME=baby-fans
PORT?=18081
IMAGE_NAME=baby-fans-backend
CONTAINER_NAME=baby-fans
DB_CONTAINER_NAME=baby-fans-db
APP_DIR=/app

# Default target
all: build

# Build the application
build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o $(BINARY_NAME) ./backend/cmd/server

# Docker compose up (with MySQL)
docker-compose-up:
	docker-compose up -d

# Docker compose down (all services)
docker-compose-down:
	docker-compose down

# Restart backend only (keeps database running)
backend-restart:
	docker-compose restart backend

# Start backend only (keeps database running)
backend-start:
	docker-compose up -d --no-deps backend

# Stop backend only
backend-stop:
	docker-compose stop backend && docker-compose rm backend

# Nginx commands
nginx-up:
	docker-compose up -d nginx

nginx-down:
	docker-compose stop nginx && docker-compose rm nginx

nginx-logs:
	docker logs -f baby-fans-nginx

nginx-restart:
	docker-compose restart nginx
