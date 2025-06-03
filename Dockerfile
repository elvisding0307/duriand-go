# 构建阶段 - 使用指定版本的Go镜像
FROM golang:1.24.3-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制项目文件
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o duriand ./cmd/main.go

# 运行阶段
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=builder /app/duriand .

# 复制配置文件
COPY .env .

# 设置GIN运行模式为release
ENV GIN_MODE=release

# 暴露端口
EXPOSE 7224

# 启动应用
CMD ["./duriand"]