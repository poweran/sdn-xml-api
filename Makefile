.PHONY: create-bin-dir install-migrate build-migrate create-docker-network start-db-server create-db migrate-up start-app

POSTGRES_PASSWORD=postgres
POSTGRES_DB=mydb
DB_CONNECTION_STRING ?= "postgres://postgres:$(POSTGRES_PASSWORD)@localhost:65432/$(POSTGRES_DB)?sslmode=disable"
MIGRATE_VERSION=latest
POSTGRES_DRIVER_VERSION=latest

create-bin-dir:
	mkdir -p bin

install-migrate:
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)
	go get -d -tags 'postgres' github.com/golang-migrate/migrate/v4/database/postgres@$(POSTGRES_DRIVER_VERSION)

build-migrate:
	go build -tags 'postgres' -ldflags="-X github.com/golang-migrate/migrate/v4/cmd/migrate.Version=$(MIGRATE_VERSION)" -o ./bin/migrate github.com/golang-migrate/migrate/v4/cmd/migrate
	chmod +x ./bin/migrate

create-docker-network:
	docker network create backend || true

start-db-server:
	docker-compose -p sdn-xml-api up -d db
	POSTGRES_DB=$(POSTGRES_DB) POSTGRES_USER=$(POSTGRES_USER) ./bin/wait-for-postgres.sh

create-db:
	PGPASSWORD=$(POSTGRES_PASSWORD) createdb -U postgres $(POSTGRES_DB) -h localhost -p 65432 || true

migrate-up:
	./bin/migrate -database $(DB_CONNECTION_STRING) -path internal/database/migrations up

migrate-down:
	./bin/migrate -database $(DB_CONNECTION_STRING) -path internal/database/migrations down

migrate-force:
	./bin/migrate -database $(DB_CONNECTION_STRING) -path internal/database/migrations force 20220420120000

start-app:
	docker-compose -p sdn-xml-api up -d app

install-all: create-bin-dir install-migrate build-migrate create-docker-network start-db-server create-db migrate-up start-app
