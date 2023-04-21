### Структура проекта:

    .
    ├── bin/
    │   └── migrate
    ├── cmd/
    │   └── myserver/
    │       └── main.go
    │       └── router.go
    ├── config/
    │   └── config.go
    │   └── config.yaml
    │   └── config-local.yaml
    ├── internal/
    │   ├── database/
    │   │   ├── migrations/
    │   │   │   └── 20220420120000_create_people_table.down.sql
    │   │   │   └── 20220420120000_create_people_table.up.sql
    │   │   ├── postgres/
    │   │   │   ├── postgres.go
    │   │   │   └── postgres_test.go
    │   │   └── repository/
    │   │       ├── people.go
    │   │       └── people_test.go
    │   └── util/
    │       ├── response_writer.go
    └── .gitignore
    └── ap.Dockerfile
    └── db.Dockerfile
    └── docker-compose.yml
    └── go.mod
    └── go.sum
    └── init-user-db.sh
    └── LICENSE
    └── Makefile
    └── README.md
    └── wait-for-postgres.sh
