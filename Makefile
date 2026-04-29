.PHONY: build test lint dev stop clean backend-build backend-test frontend-build frontend-test docker-build docker-up docker-down

# Backend
backend-build:
	cd backend && go build -o ../bin/ems-api ./main.go

backend-test:
	cd backend && go test ./... -v

backend-lint:
	cd backend && go vet ./...

# Frontend
frontend-build:
	cd frontend && npm run build

frontend-test:
	cd frontend && npm test

frontend-lint:
	cd frontend && npx vue-tsc --noEmit

# Combined
build: backend-build frontend-build
test: backend-test frontend-test
lint: backend-lint frontend-lint

# Development
dev:
	./start-dev.sh

stop:
	./stop-dev.sh

clean:
	rm -rf bin/
	cd frontend && rm -rf dist/

# Docker
docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-restart: docker-down docker-up

docker-logs:
	docker compose logs -f

docker-clean:
	docker compose down -v
