FROM golang:1.23.4-alpine3.19

# Рабочая директория
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем приложение
RUN go build -o main ./main.go
COPY envs/.env ./envs/.env

# Открываем порт приложения
EXPOSE 3001

CMD ["./main"]

