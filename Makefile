.PHONY: compose up

compose up:
	docker compose --env-file ./configs/envs/local.env --file ./deployments/docker-compose.yml -p merch-shop up --force-recreate --build --no-deps