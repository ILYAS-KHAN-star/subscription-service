.PHONY: run build clean test swagger migrate-up migrate-down docker-up docker-down deps help

help:
 @echo "Available commands:"
 @echo "  make run          - Run application locally"
 @echo "  make build        - Build binary"
 @echo "  make clean        - Clean build artifacts"
 @echo "  make test         - Run tests"
 @echo "  make swagger      - Generate swagger documentation"
 @echo "  make migrate-up   - Run database migrations up"
 @echo "  make migrate-down - Run database migrations down"
 @echo "  make docker-up    - Start services with Docker Compose"
 @echo "  make docker-down  - Stop Docker services"
 @echo "  make deps         - Install dependencies"

run:
 go run cmd/app/main.go

build:
 go build -o bin/app cmd/app/main.go

clean:
 rm -rf bin/
 go clean

test:
 go test -v -cover ./...

swagger:
 swag init -g cmd/app/main.go -o docs

migrate-up:
 migrate -path migrations -database "postgresql://postgres:password@localhost:5432/subscriptions?sslmode=disable" up

migrate-down:
 migrate -path migrations -database "postgresql://postgres:password@localhost:5432/subscriptions?sslmode=disable" down

docker-up:
 docker-compose up -d

docker-down:
 docker-compose down

docker-build:
 docker-compose build

docker-logs:
 docker-compose logs -f app

deps:
 go mod download
 go mod tidy
 go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
 go install github.com/swaggo/swag/cmd/swag@latest