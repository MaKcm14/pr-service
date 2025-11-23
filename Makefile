build:
	docker compose -f docker-compose-db.yml build
	docker compose -f docker-compose-pr-service.yml build
	docker compose -f docker-compose-migrations.yml build

start-db:
	docker compose -f docker-compose-db.yml up -d

migrate:
	docker compose -f docker-compose-migrations.yml up -d

start-service:
	docker compose -f docker-compose-pr-service.yml up -d

stop:
	docker compose -f docker-compose-pr-service.yml down
	docker compose -f docker-compose-db.yml down
