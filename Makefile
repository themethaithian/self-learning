.PHONY: up test vet run

up:
	docker compose up -d

test:
	go test ./...

vet:
	go vet ./...

run:
	go run ./cmd/api
