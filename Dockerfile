# 构建阶段
FROM golang:1.25.3-alpine3.22 AS builder

# 安装编译依赖
RUN apk add --no-cache gcc g++ cmake git make unzip

# 设置工作目录
WORKDIR /gpress

# 设置国内代理
#RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 复制项目代码
COPY . .

# 编译项目
RUN go build --tags "fts5" -ldflags "-w -s" -o gpress

# 初始化文件
RUN rm -rf /gpress/gpressdatadir/dict && \
    unzip /gpress/gpressdatadir/dict.zip -d /gpress/gpressdatadir && \
    rm -rf /gpress/gpressdatadir/dict.zip && \
    cp -rf /gpress/gpressdatadir/fts5/libsimple.so /gpress/gpressdatadir/ && \
    rm -rf /gpress/gpressdatadir/fts5 && \
    mkdir -p /gpress/gpressdatadir/fts5 && \
    mv /gpress/gpressdatadir/libsimple.so /gpress/gpressdatadir/fts5/


# 运行阶段
FROM alpine:3.22.2

# 安装运行时依赖
RUN apk add --no-cache libgcc libstdc++ sqlite-libs


# 设置工作目录
WORKDIR /gpress

RUN mkdir -p ./gpressdatadir

# 复制编译产物
COPY --from=builder /gpress/gpress .
COPY --from=builder /gpress/gpressdatadir ./gpressdatadir/



# 暴露端口
EXPOSE 660

# 设置数据卷
VOLUME ["/gpress/gpressdatadir"]

# 启动命令
CMD ["./gpress"]