# Build stage
FROM golang:1.24.4-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/price-producer ./cmd/price-producer


FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache netcat-openbsd

COPY --from=builder /app/bin/price-producer /app/price-producer
COPY config/ /app/config/

EXPOSE 8080

ENTRYPOINT ["/app/price-producer"]




