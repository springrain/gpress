# 构建阶段
FROM golang:1.25.3-alpine3.22 AS builder

# 安装编译依赖
RUN apk add --no-cache gcc g++ cmake git make zip unzip

# 设置工作目录
WORKDIR /app

# 设置国内代理
#RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 复制项目代码
COPY . .

# 编译项目
RUN go build --tags "fts5" -ldflags "-w -s" -o gpress

# 初始化文件
RUN rm -rf /app/gpressdatadir/dict && \
    unzip /app/gpressdatadir/dict.zip -d /app/gpressdatadir && \
    rm -rf /app/gpressdatadir/dict.zip && \
    cp -rf /app/gpressdatadir/fts5/libsimple.so /app/gpressdatadir/ && \
    rm -rf /app/gpressdatadir/fts5 && \
    mkdir -p /app/gpressdatadir/fts5 && \
    mv /app/gpressdatadir/libsimple.so /app/gpressdatadir/fts5/


# 运行阶段
FROM alpine:3.22.2

# 安装运行时依赖
RUN apk add --no-cache libgcc libstdc++ sqlite-libs


# 设置工作目录
WORKDIR /app

RUN mkdir -p ./gpressdatadir

# 复制编译产物
COPY --from=builder /app/gpress .
COPY --from=builder /app/gpressdatadir ./gpressdatadir/



# 暴露端口
EXPOSE 660

# 设置数据卷
VOLUME ["/app/gpressdatadir"]

# 启动命令
CMD ["./gpress"]