# Dockerfile для базы данных Postgres
FROM postgres:14-alpine

# переменная среды для установки пароля Postgres
ENV POSTGRES_PASSWORD postgres
