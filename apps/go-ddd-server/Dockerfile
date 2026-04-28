FROM golang:1.25.1 AS builder

WORKDIR /app

COPY . .

WORKDIR /app/apps/go-ddd-server

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/apps/go-ddd-server/main .

CMD ["./main"]
