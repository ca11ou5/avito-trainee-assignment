.PHONY: compose up

compose-up:
	docker compose --env-file ./configs/envs/local.env --file ./deployments/docker-compose.yml -p merch-shop up --force-recreate --build --no-deps -d

lint:
	golangci-lint run

tests:
	go test -short -coverprofile=coverage.out ./...

testscover: tests
	go tool cover -func=coverage.out

e2e-tests: compose-up
	sleep 30
	go test -v ./e2e/
	docker compose --env-file ./configs/envs/local.env --file ./deployments/docker-compose.yml -p merch-shop down
