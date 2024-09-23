up-dev:
	docker compose -f docker-compose.dev.yml --build up

up-prod:
	docker-compose -f docker-compose.prod.yml --build up

down:
	docker compose -f docker-compose.dev.yml down