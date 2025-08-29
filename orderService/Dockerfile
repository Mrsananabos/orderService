FROM golang:1.24.2-bookworm AS builder
    WORKDIR /app

    COPY ./orderService/ ./

    RUN go build -o orderService ./cmd/app/main.go

FROM debian:bookworm

    WORKDIR /app
    RUN apt update && apt install make
    COPY --from=builder /app/orderService .
    COPY --from=builder /app/migrations ./migrations

    CMD ["./orderService"]

