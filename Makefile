build:
	docker compose -f docker-compose.yml build
	docker compose -f docker-compose-migrations.yml build

start:
	docker compose -f docker-compose.yml up -d

migrate:
	docker compose -f docker-compose-migrations.yml up -d

stop:
	docker compose -f docker-compose.yml down
