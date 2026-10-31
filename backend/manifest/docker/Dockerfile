# 构建阶段
FROM golang:1.24-alpine AS builder

WORKDIR /app
# 先复制依赖文件，充分利用 Docker 层缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并执行静态交叉编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

# 运行阶段
FROM alpine:3.19

# 安装基础证书与时区数据，确保能安全发起外网 HTTPS 请求（Gemini / OpenAI / 豆包）
RUN apk --no-cache add ca-certificates tzdata curl wget \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app
COPY --from=builder /app/main .
# 复制配置文件模板与 ComfyUI 工作流
COPY manifest/config/config*.yaml manifest/config/
COPY workflows workflows/
# 若本地未挂载自定义 config.yaml，则以模板兜底初始化
RUN [ -f manifest/config/config.yaml ] || cp manifest/config/config.example.yaml manifest/config/config.yaml

EXPOSE 8000

# 容器健康检查
HEALTHCHECK --interval=15s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:8000/hello || exit 1

CMD ["./main"]
