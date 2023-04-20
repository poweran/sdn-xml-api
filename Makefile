.PHONY: install create-bin-dir install-migrate build-migrate create-db migrate-up

PG_PASSWORD=postgres
DB_CONNECTION_STRING ?= "postgres://postgres:$(PG_PASSWORD)@localhost:65432/mydb?sslmode=disable"
MIGRATE_VERSION=latest
POSTGRES_DRIVER_VERSION=latest

install:
	go install github.com/example/myapp

create-bin-dir:
	mkdir -p bin

install-migrate:
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)
	go get -d -tags 'postgres' github.com/golang-migrate/migrate/v4/database/postgres@$(POSTGRES_DRIVER_VERSION)

build-migrate:
	go build -tags 'postgres' -ldflags="-X github.com/golang-migrate/migrate/v4/cmd/migrate.Version=$(MIGRATE_VERSION)" -o ./bin/migrate github.com/golang-migrate/migrate/v4/cmd/migrate
	chmod +x ./bin/migrate

create-db:
	PGPASSWORD=$(PG_PASSWORD) createdb -U postgres mydb -h localhost -p 65432 || true

migrate-up:
	./bin/migrate -database $(DB_CONNECTION_STRING) -path internal/database/migrations up

migrate-down:
	./bin/migrate -database $(DB_CONNECTION_STRING) -path internal/database/migrations down

migrate-force:
	./bin/migrate -database $(DB_CONNECTION_STRING) -path internal/database/migrations force 20220420120000

install-all: create-bin-dir install-migrate build-migrate create-db migrate-up
