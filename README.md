Структура проекта:
├── bin/
│   └── migrate
├── cmd/
│   └── myserver/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
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
│   ├── handler/
│   │   ├── get_names.go
│   │   ├── state.go
│   │   └── update.go
│   ├── model/
│   │   └── person.go
│   ├── service/
│   │   ├── ofac/
│   │   │   ├── ofac.go
│   │   │   └── ofac_test.go
│   │   └── person/
│   │       ├── person.go
│   │       └── person_test.go
│   └── util/
│       ├── json.go
│       ├── response.go
│       └── util_test.go
├── static/
│   └── index.html
└── go.mod