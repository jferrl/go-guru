FROM golang:1.22-alpine as builder

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 go build -o go-guru ./cmd/go-guru

FROM alpine:latest

COPY --from=builder /app/go-guru /go-guru

RUN chmod +x ./go-guru

ENTRYPOINT ["./go-guru"]