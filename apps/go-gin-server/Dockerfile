FROM golang:1.25.1 AS builder

WORKDIR /app

# 複製整個 monorepo
COPY . .

WORKDIR /app/apps/go-counter-server

# 下載依賴
RUN go mod download

# 編譯執行檔
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# 使用更小的 base image
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata openssl

WORKDIR /app

# 只複製執行檔
COPY --from=builder /app/apps/go-counter-server/main .


CMD ["./main"]