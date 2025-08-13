
FROM golang:1.24.4-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api-gateway ./cmd/api-gateway

FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache netcat-openbsd

COPY --from=builder /app/bin/api-gateway /app/api-gateway
COPY config/ /app/config/

EXPOSE 8081

ENTRYPOINT ["/app/api-gateway"]




