FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN go mod download

# Собираем основной бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd

# Собираем миграции
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o migrate ./cmd/migrations

# ---

FROM alpine:latest

WORKDIR /root/

# Устанавливаем bash и make (если используешь Makefile)
RUN apk add --no-cache bash

# Копируем исполняемые файлы
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY .env .
COPY ./migrations ./migrations
COPY wait-for-it.sh .

RUN chmod +x ./main ./migrate ./wait-for-it.sh

# Запускаем: ждем БД, выполняем миграцию и стартуем приложение
CMD ["sh", "-c", "./wait-for-it.sh postgres:5432 --strict --timeout=30 -- ./migrate -up && ./main"]