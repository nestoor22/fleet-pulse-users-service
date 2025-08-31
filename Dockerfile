# Build stage
FROM public.ecr.aws/docker/library/golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Install goose CLI
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy source and build your app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s' \
    -o main ./cmd/server/main.go

# Runtime stage
FROM public.ecr.aws/docker/library/alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN mkdir -p /app && chown appuser:appgroup /app

WORKDIR /app

# Copy app binary
COPY --from=builder /build/main /app/main
# Copy goose binary from Go bin dir
COPY --from=builder /go/bin/goose /usr/local/bin/goose

COPY migrations /app/migrations
COPY .env /app/.env

# Copy start.sh
COPY start.sh /app/start.sh
RUN chmod +x /app/start.sh /app/main && chown appuser:appgroup /app/*
RUN chown -R appuser:appgroup /app/migrations
RUN chown appuser:appgroup /app/.env

USER appuser

EXPOSE 8000

ENTRYPOINT ["/app/start.sh"]
