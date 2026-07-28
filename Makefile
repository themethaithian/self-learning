.PHONY: up test vet run dev dev-down dev-reset seed logs

up:
	docker compose up -d

test:
	go test ./...

vet:
	go vet ./...

run:
	go run ./cmd/api

dev:
	docker compose up -d --build

dev-down:
	docker compose down

dev-reset:
	docker compose down -v

seed:
	docker compose run --rm seed

logs:
	docker compose logs -f
