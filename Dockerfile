# Этап 1: Сборка приложения
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код проекта
COPY . .

# Собираем приложение. CGO_ENABLED=0 создает статически скомпилированный бинарник.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /server ./cmd/server/main.go

# Этап 2: Создание минимального финального образа
FROM alpine:latest

WORKDIR /root/

# Копируем скомпилированный бинарник из этапа сборки
COPY --from=builder /server .

# Открываем порт, на котором будет работать сервер
EXPOSE 8181

# Команда для запуска сервера
CMD ["./server"]