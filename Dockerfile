# 使用 go1.20 alphine 为命令行打包镜像
FROM golang:1.20-alpine AS builder
WORKDIR /
COPY . .
RUN go build -o cubox-archiver cmd/main.go


FROM alpine:latest
# 安装 CA 证书，用于 HTTPS 请求
RUN apk update && apk add --no-cache ca-certificates

# 将二进制文件从 build 阶段复制到生产镜像中
COPY --from=builder /cubox-archiver /cubox-archiver

# 容器启动时运行的命令
ENTRYPOINT ["./cubox-archiver"]
