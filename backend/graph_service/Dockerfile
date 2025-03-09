FROM golang:1.23.4-alpine3.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY graph_service/ ./graph_service

WORKDIR /app/graph_service
RUN go build -o main ./main.go

COPY graph_service/.env ./

EXPOSE 3001
CMD ["./main"]


