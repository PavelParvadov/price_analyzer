
FROM golang:1.24.4-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/price-cache ./cmd/price-cache


FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache netcat-openbsd

COPY --from=builder /app/bin/price-cache /app/price-cache
COPY config/ /app/config/

EXPOSE 6000

ENTRYPOINT ["/app/price-cache"]




