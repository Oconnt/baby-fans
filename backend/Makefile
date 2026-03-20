.PHONY: build run dev test clean docker-build docker-run docker-stop docker-logs

# Build variables
BINARY_NAME=baby-fans
PORT=18081
IMAGE_NAME=baby-fans-backend
CONTAINER_NAME=baby-fans

# Default target
all: build

# Build the application
build:
	CGO_ENABLED=0 GOOS=linux go build -o $(BINARY_NAME) ./cmd/server

# Run locally (requires SQLite database)
run: build
	./$(BINARY_NAME)

# Run in development mode
dev:
	go run ./cmd/server

# Run tests
test:
	go test -v ./...

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
		-v $(PWD)/storage:/app/storage \
		-v $(PWD)/config:/app/config \
		-v $(PWD)/baby-fans.db:/app/baby-fans.db \
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
