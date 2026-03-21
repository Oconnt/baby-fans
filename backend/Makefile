.PHONY: build run dev test clean docker-build docker-build-linux docker-run docker-run-mysql docker-stop docker-logs docker-compose-up docker-compose-down docker-shell docker-push migrate

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

# Run database migration (requires local Go installation)
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

# Run Docker container with MySQL (requires mysql container running)
docker-run-mysql:
	@docker stop $(DB_CONTAINER_NAME) 2>/dev/null || true
	@docker rm $(DB_CONTAINER_NAME) 2>/dev/null || true
	@docker network rm baby-fans-net 2>/dev/null || true
	docker network create baby-fans-net
	docker run -d --name $(DB_CONTAINER_NAME) \
		--network baby-fans-net \
		-e MYSQL_ROOT_PASSWORD=baby_fans_password \
		-e MYSQL_DATABASE=baby_fans \
		-v mysql_data:/var/lib/mysql \
		--restart unless-stopped \
		mysql:8.0
	@echo "Waiting for MySQL to be ready..."
	@for i in $$(seq 1 30); do \
		docker exec $(DB_CONTAINER_NAME) mysqladmin ping -h localhost -u root -pbaby_fans_password >/dev/null 2>&1 && break || sleep 1; \
	done
	@docker exec $(DB_CONTAINER_NAME) mysqladmin ping -h localhost -u root -pbaby_fans_password >/dev/null 2>&1 || (echo "MySQL failed to start" && exit 1)
	@echo "MySQL is ready!"
	docker run -d --name $(CONTAINER_NAME) \
		--network baby-fans-net \
		-p $(PORT):18081 \
		-v $(shell pwd)/storage:$(APP_DIR)/storage \
		--restart unless-stopped \
		$(IMAGE_NAME):latest

# Run Docker container (SQLite fallback mode)
docker-run:
	docker run -d --name $(CONTAINER_NAME) \
		-p $(PORT):18081 \
		-v $(shell pwd)/storage:$(APP_DIR)/storage \
		--restart unless-stopped \
		$(IMAGE_NAME):latest

# Stop Docker containers
docker-stop:
	docker stop $(CONTAINER_NAME) 2>/dev/null || true
	docker rm $(CONTAINER_NAME) 2>/dev/null || true
	docker stop $(DB_CONTAINER_NAME) 2>/dev/null || true
	docker rm $(DB_CONTAINER_NAME) 2>/dev/null || true
	docker network rm baby-fans-net 2>/dev/null || true

# View backend logs
docker-logs:
	docker logs -f $(CONTAINER_NAME)

# Docker compose up (with MySQL)
docker-compose-up:
	docker-compose up -d

# Docker compose down
docker-compose-down:
	docker-compose down

# Shell into backend container
docker-shell:
	docker exec -it $(CONTAINER_NAME) /bin/sh

# Shell into database container
docker-db-shell:
	docker exec -it $(DB_CONTAINER_NAME) mysql -u root -pbaby_fans_password baby_fans

# Push Docker image to registry
docker-push:
	docker push $(IMAGE_NAME):latest
