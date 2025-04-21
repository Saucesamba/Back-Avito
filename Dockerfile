# Используем базовый образ с Go (укажи нужную версию)
FROM golang:1.23-alpine AS builder

# Устанавливаем зависимости
RUN apk update && apk add --no-cache git

# Устанавливаем рабочий каталог внутри контейнера
WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/main/main.go

# Создаем финальный образ
FROM alpine:latest

# Копируем исполняемый файл из builder-образа
COPY --from=builder /app/main /app/main

# Устанавливаем рабочий каталог
WORKDIR /app

# Копируем конфигурационный файл
COPY config.yaml config.yaml

#Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["/app/main"]