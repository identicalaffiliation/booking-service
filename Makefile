.PHONY: db_up
db_up:
	docker compose up --build -d booking-postgresql
	docker compose up --build -d notifications-postgresql

.PHONY: db_down
db_down:
	docker compose down booking-postgresql
	docker compose down notifications-postgresql

.PHONY: clean
clean:
	docker compose down -v

.PHONY: migrate_up
migrate_up:
	docker compose run --build --rm booking-migrator ./migrator -a up

.PHONY: migrate_down
migrate_down:
	docker compose run --build --rm booking-migrator ./migrator -a down

.PHONY: migrate_reset
migrate_reset:
	docker compose run --build --rm booking-migrator ./migrator -a reset

.PHONY: up
up:
	docker compose up --build -d booking-api