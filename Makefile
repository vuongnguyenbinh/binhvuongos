.PHONY: build run docker-build docker-up docker-down

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
