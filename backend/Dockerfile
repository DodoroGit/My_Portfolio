# 使用官方 Golang 基礎映像
FROM golang:1.22-alpine

# 安裝必要套件
RUN apk --no-cache add ca-certificates

# 設定工作目錄
WORKDIR /app

# 複製 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 複製 backend 的程式碼
COPY . .

# 編譯 Golang 程式
RUN go build -o backend main.go

# 設定環境變數
ENV PORT=8080

# 開放 port
EXPOSE 8080

# 啟動
CMD ["./backend"]
