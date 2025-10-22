FROM golang:1.23-alpine

RUN apk add --no-cache wget

WORKDIR /app

COPY server.go .

RUN go build -o mock-server server.go

EXPOSE 8080

CMD ["./mock-server"]
