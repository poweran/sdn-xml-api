FROM golang:1.20-alpine

# Устанавливаем зависимости
RUN apk update && apk add --no-cache make git postgresql-client

# Копируем все файлы приложения внутрь контейнера
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Install all dependencies
RUN make install-all

CMD ["bin/myserver"]
