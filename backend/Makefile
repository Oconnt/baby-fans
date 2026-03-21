.PHONY: build run dev test clean docker-build docker-build-linux docker-run docker-stop docker-logs docker-compose-up docker-compose-down docker-shell docker-push migrate

# Build variables
BINARY_NAME=baby-fans
PORT?=18081
IMAGE_NAME=baby-fans-backend
CONTAINER_NAME=baby-fans
APP_DIR=/app

# Default target
all: build

# Build the application
build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o $(BINARY_NAME) ./cmd/server

# Run locally (requires SQLite database)
run: build
	./$(BINARY_NAME)

# Run in development mode
dev:
	go run ./cmd/server

# Run tests
test:
	go test -v ./...

# Run database migration
migrate:
	go run ./scripts/migrate.go

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f *.db

# Build Docker image
docker-build:
	docker build -t $(IMAGE_NAME):latest .

docker-build-linux:
	docker build --platform=linux/amd64 -t $(IMAGE_NAME):latest .

# Run Docker container
docker-run:
	docker run -d --name $(CONTAINER_NAME) \
		-p $(PORT):18081 \
		-v $(shell pwd)/storage:$(APP_DIR)/storage \
		-v $(shell pwd)/etc:$(APP_DIR)/config \
		-v $(shell pwd)/baby-fans.db:$(APP_DIR)/baby-fans.db \
		-e SERVER_PORT=18081 \
		-e SERVER_MODE=release \
		--restart unless-stopped \
		$(IMAGE_NAME):latest

# Stop Docker container
docker-stop:
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

# View Docker logs
docker-logs:
	docker logs -f $(CONTAINER_NAME)

# Docker compose up
docker-compose-up:
	docker-compose up -d

# Docker compose down
docker-compose-down:
	docker-compose down

# Shell into Docker container
docker-shell:
	docker exec -it $(CONTAINER_NAME) /bin/sh

# Push Docker image to registry
docker-push:
	docker push $(IMAGE_NAME):latest
