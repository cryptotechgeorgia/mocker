# Project variables
APP_NAME := mocker
DOCKER_COMPOSE_FILE := docker-compose.build.yml
DOCKER_REGISTRY := dcreg.coinet.ge
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Targets
.PHONY: all build run test clean docker-build docker-up docker-down format vet lint

build:
	@echo "Building the Go project..."
	@go build -o $(APP_NAME) main.go

run: build
	@echo "Running the Go project..."
	@./$(APP_NAME)

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@echo "Cleaning up..."
	@go clean
	@rm -f $(APP_NAME)

deploy:
	@echo "Building Docker image..."
	@read -p "Enter Docker image tag : " IMAGE_TAG; \
	docker build -t $(DOCKER_REGISTRY)/$(APP_NAME):$$IMAGE_TAG . ;\
       	docker push  $(DOCKER_REGISTRY)/$(APP_NAME):$$IMAGE_TAG


docker-up:
	@echo "Starting Docker containers..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down

format:
	@echo "Formatting Go files..."
	@gofmt -w $(GO_FILES)

vet:
	@echo "Running go vet..."
	@go vet ./...

lint:
	@echo "Running linter..."
	@golangci-lint run ./...

all: clean build run

