# ==========================================
# 第一阶段：构建阶段（Builder）
# 使用带有 Go 环境的镜像来编译我们的代码
# ==========================================
FROM golang:alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装 git（如果需要从私有仓库拉取依赖）和其他必要的工具
RUN apk add --no-cache git
# 配置 Go 代理（国内加速下载依赖）
ENV GOPROXY=https://proxy.golang.org,direct

# 先复制 go.mod 和 go.sum 并下载依赖（利用 Docker 缓存加速）
COPY go.mod go.sum ./
RUN go mod download

# 把项目所有源代码复制进容器
COPY . .

# 编译 Go 语言程序，生成一个名为 main 的可执行文件
RUN go build -o main ./cmd


# ==========================================
# 第二阶段：运行阶段（Runner）
# 使用一个极小的 alpine 镜像来运行编译好的文件，减小最终体积
# ==========================================
FROM alpine:latest

WORKDIR /app

# 从构建阶段把编译好的二进制文件复制过来
COPY --from=builder /app/main .
# 复制配置文件和 Swagger 文档（运行程序必须的静态文件）
COPY --from=builder /app/config ./config
COPY --from=builder /app/docs ./docs

# 声明容器内部监听的端口
EXPOSE 8080

# 启动容器时执行的命令
CMD ["./main"]