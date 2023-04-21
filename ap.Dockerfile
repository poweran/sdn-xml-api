# базовый образ
FROM golang:1.20-alpine

# установка зависимостей
RUN apk update && apk add --no-cache git make

# создание рабочей директории
WORKDIR /app

# копирование исходного кода приложения в рабочую директорию
COPY . .

# установка зависимостей Go
RUN go mod download

# сборка приложения
RUN go build -o myserver sdn-xml-api/cmd/myserver

#
EXPOSE 8080

# запуск приложения
CMD [ "/app/myserver" ]
