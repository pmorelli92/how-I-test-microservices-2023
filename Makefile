include local.env
export

dependencies:
	docker compose up -d

test:
	go test -tags=accept --count=1 ./accept_test/... -v

test/docker:
	docker compose --file docker-compose.yaml --file docker-compose-test.yaml  up --build -d
	docker attach accept_test
