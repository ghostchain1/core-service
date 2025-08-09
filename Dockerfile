# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-X github.com/ghostchain1/core-service/pkg/version.Version=0.1.0" -o core-service ./cmd/core-service

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/core-service .

EXPOSE 8080

CMD ["./core-service"]
