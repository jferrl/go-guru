FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go build -o go-guru ./cmd/go-guru

ENTRYPOINT ["./go-guru"]