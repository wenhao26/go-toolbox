FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o consumer main.go

CMD ["./consumer", "--config=config.yaml"]