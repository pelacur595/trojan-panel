FROM golang:1.18.3-alpine3.15 AS builder
RUN mkdir /tpdata/trojan-panel/app/
ADD . /tpdata/trojan-panel/app/
WORKDIR /tpdata/trojan-panel/app/
RUN go install mvdan.cc/garble@latest SET CGO_ENABLED=0 SET GOOS=linux SET GOARCH=amd64 garble -literals build -o build/trojan-panel

FROM alpine:3.15
RUN mkdir /tpdata/trojan-panel/app/
LABEL maintainer="jonsosnyan <https://jonssonyan.com>"
ENV mariadb_ip=trojan-panel-mariadb \
    mariadb_port=3306 \
    mariadb_user=root \
    mariadb_pas=123456 \
    redis_host=my-redis \
    redis_port=6379 \
    redis_pass=123456
WORKDIR mkdir /tpdata/trojan-panel/app/
ADD --from=builder build/trojan-panel .
RUN apk add bash tzdata ca-certificates && \
    rm -rf /var/cache/apk/*
ENTRYPOINT ./trojan-panel \
    -host=${mariadb_ip} \
    -port=${mariadb_port} \
    -user=${mariadb_user} \
    -password=${mariadb_pas} \
    -redisHost=${redis_host} \
    -redisPort=${redis_port} \
    -redisPassword=${redis_pass}