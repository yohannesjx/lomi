.PHONY: run build docker-up docker-down

run:
	cd backend && go run cmd/api/main.go

build:
	cd backend && go build -o bin/api cmd/api/main.go

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

tidy:
	cd backend && go mod tidy
