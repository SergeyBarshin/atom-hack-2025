# Фаза сборки
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем весь проект
COPY . .
RUN go mod tidy
# Переходим в директорию сервиса и собираем
WORKDIR /app/auth-service
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

# Фаза запуска
FROM alpine:latest AS runner

WORKDIR /app

# Копируем бинарник и .env файл
COPY --from=builder /app/auth-service/main .
COPY --from=builder /app/.env .

# Устанавливаем пакет для работы с .env
RUN apk add --no-cache libc6-compat

EXPOSE 3002
CMD ["./main"]