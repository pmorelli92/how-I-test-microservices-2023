include local.env
export

dependencies:
	docker compose up -d

run:
	go run cmd/main.go

tests:
	go test -tags=accept --count=1 ./...
