FROM golang:1.17

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY src src

RUN go build -o main src/main.go

EXPOSE 8080

CMD ["./main"]
