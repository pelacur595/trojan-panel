FROM alpine
LABEL maintainer="jonsosnyan <https://jonssonyan.com>"

ENV mariadb_ip=trojan-panel-mariadb \
    mariadb_port=3306 \
    mariadb_user=root \
    mariadb_pas=123456 \
    redis_host=my-redis \
    redis_port=6379 \
    redis_pass=123456

WORKDIR /tpdata/trojan-panel/
ADD build/trojan-panel .
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