FROM alpine:3.15
LABEL maintainer="jonsosnyan <https://jonssonyan.com>"
RUN mkdir -p /tpdata/trojan-panel/
WORKDIR /tpdata/trojan-panel/
ENV mariadb_ip=127.0.0.1 \
    mariadb_port=9507 \
    mariadb_user=root \
    mariadb_pas=123456 \
    redis_host=127.0.0.1 \
    redis_port=6378 \
    redis_pass=123456 \
    TZ=Asia/Shanghai
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
COPY build/trojan-panel-${TARGETOS}-${TARGETARCH}${TARGETVARIANT} trojan-panel
# 国内环境开启以下注释 设置apk国内镜像
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add bash tzdata ca-certificates && \
    rm -rf /var/cache/apk/*
ENTRYPOINT chmod 777 ./trojan-panel && \
    ./trojan-panel \
    -host=${mariadb_ip} \
    -port=${mariadb_port} \
    -user=${mariadb_user} \
    -password=${mariadb_pas} \
    -redisHost=${redis_host} \
    -redisPort=${redis_port} \
    -redisPassword=${redis_pass}