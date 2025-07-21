include .env
export

env:
	cp .env.example .env

build:
	docker compose up --build -d && docker compose logs -f

up:
	docker compose up -d && docker compose logs -f

stop:
	docker compose stop

lint:
	golangci-lint run

fmt:
	golangci-lint fmt

swag:
	swag init -g cmd/app/main.go