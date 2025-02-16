FROM golang:1.23 AS builder

WORKDIR /app
COPY . .

RUN go mod download

RUN go build -o ./build/merch-store ./cmd/merch-store \
    && go clean -cache -modcache

# FROM alpine

# WORKDIR /

# COPY --from=builder /build/merch-store /merch-store

EXPOSE 8080

CMD ["./build/merch-store"]