FROM golang:1.25.1 AS builder

WORKDIR /app

# Copy entire monorepo
COPY . .

WORKDIR /app/apps/go-layered-server

# Download dependencies
RUN go mod download

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Use smaller base image
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata openssl

WORKDIR /app

# Copy only the binary
COPY --from=builder /app/apps/go-layered-server/main .


CMD ["./main"]