# ビルドステージ
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 依存関係のコピーとダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピーとビルド
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 実行ステージ
FROM alpine:3.19

WORKDIR /app

# ビルドしたバイナリをコピー
COPY --from=builder /app/main .

# 実行時の設定
EXPOSE 8080
CMD ["./main"] 