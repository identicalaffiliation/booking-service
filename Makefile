.PHONY: db_up
db_up:
	docker compose up --build -d booking-postgresql
	docker compose up --build -d notifications-postgresql

.PHONY: db_down
db_down:
	docker compose down booking-postgresql
	docker compose down notifications-postgresql

.PHONY: kafka_up
kafka_up:
	docker compose up --build -d zookeeper
	docker compose up --build -d kafka
	docker compose up --build -d kafka-ui

.PHONY: clean
clean:
	docker compose down -v

.PHONY: migrate_up
migrate_up:
	docker compose run --build --rm booking-migrator ./migrator -a up
	docker compose run --build --rm notifications-migrator ./migrator -a up

.PHONY: migrate_down
migrate_down:
	docker compose run --build --rm booking-migrator ./migrator -a down
	docker compose run --build --rm notifications-migrator ./migrator -a down

.PHONY: migrate_reset
migrate_reset:
	docker compose run --build --rm booking-migrator ./migrator -a reset
	docker compose run --build --rm notifications-migrator ./migrator -a reset
.PHONY: up
up:
	docker compose up --build -d booking-api
	docker compose up --build -d notifications-service
.PHONY: down
down:
	docker compose down booking-api
	docker compose down notifications-service

.PHONY: integration_tests
integration_tests:
	cd booking && go test ./tests/integration/...

.PHONY: behavioural_tests
behavioural_tests:
	cd booking/tests/behavioural && go test ./...