# Dockerfile для базы данных Postgres
FROM postgres:14-alpine

# добавление скрипта для создания базы данных
COPY ./init-user-db.sh /docker-entrypoint-initdb.d/init-user-db.sh

# переменная среды для установки пароля Postgres
ENV POSTGRES_PASSWORD postgres
