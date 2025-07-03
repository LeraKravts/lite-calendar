FROM golang:1.24-bullseye

# Установим goose прямо в контейнер
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Создаём директорию под миграции
WORKDIR /app


ENTRYPOINT ["/go/bin/goose"]
