FROM golang:1.23.4-alpine3.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY auth-service/ ./auth-service

WORKDIR /app/auth-service
RUN go build -o main ./main.go

COPY auth-service/.env ./

EXPOSE 3002
CMD ["./main"]

