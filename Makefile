.PHONY: build run docker-build docker-up docker-down migrate-up migrate-down sqlc

# Build in Docker (no local Go needed)
docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

# Full deploy cycle
deploy: docker-build docker-up
	@echo "✓ Deployed at http://localhost:3000"

# Database migrations (requires migrate CLI and DATABASE_URL)
migrate-up:
	migrate -database "$$DATABASE_URL" -path internal/db/migrations up

migrate-down:
	migrate -database "$$DATABASE_URL" -path internal/db/migrations down 1

# Generate sqlc code (requires sqlc CLI)
sqlc:
	cd internal/db && sqlc generate
