.PHONY: help build run stop clean test lint

help:
    @echo "Available commands:"
    @echo "  make build      - Build Docker images"
    @echo "  make up         - Start containers"
    @echo "  make down       - Stop containers"
    @echo "  make logs       - View container logs"
    @echo "  make clean      - Remove all containers and volumes"
    @echo "  make test       - Run tests"
    @echo "  make lint       - Run linter"
    @echo "  make seed-env   - Create .env from .env.example"

build:
    docker-compose build

up:
    docker-compose up -d

down:
    docker-compose down

logs:
    docker-compose logs -f app

clean:
    docker-compose down -v

test:
    go test ./...

lint:
    golangci-lint run ./...

seed-env:
    cp .env.example .env
    @echo ".env file created. Update values as needed."

fmt:
    go fmt ./...

mod-tidy:
    go mod tidy