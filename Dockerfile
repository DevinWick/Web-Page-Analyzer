# Build stage
FROM golang:1.24.4-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64


WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-w -s" -o web-analyzer

# Last stage
FROM alpine:latest

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/web-analyzer ./
COPY --from=builder /app/pages ./pages
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./web-analyzer"]