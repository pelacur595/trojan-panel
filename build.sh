#!/usr/bin/env bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

init_var() {
  ECHO_TYPE="echo -e"

  trojan_panel_version=1.1.3

  arch_arr=('amd64' 'arm64')

  touch Dockerfile
}

echo_content() {
  case $1 in
  "red")
    ${ECHO_TYPE} "\033[31m$2\033[0m"
    ;;
  "green")
    ${ECHO_TYPE} "\033[32m$2\033[0m"
    ;;
  "yellow")
    ${ECHO_TYPE} "\033[33m$2\033[0m"
    ;;
  "blue")
    ${ECHO_TYPE} "\033[34m$2\033[0m"
    ;;
  "purple")
    ${ECHO_TYPE} "\033[35m$2\033[0m"
    ;;
  "skyBlue")
    ${ECHO_TYPE} "\033[36m$2\033[0m"
    ;;
  "white")
    ${ECHO_TYPE} "\033[37m$2\033[0m"
    ;;
  esac
}

main() {
  for get_arch in ${arch_arr[*]}; do
    if [[ ! -f build/trojan-panel-linux-${get_arch} ]]; then
      continue
    fi
    echo_content skyBlue "开始构建trojan-panel-linux-${get_arch}"

    cat >Dockerfile <<EOF
FROM alpine:3.15
LABEL maintainer="jonsosnyan <https://jonssonyan.com>"
RUN mkdir -p /tpdata/trojan-panel/
WORKDIR /tpdata/trojan-panel/
ENV mariadb_ip=trojan-panel-mariadb \
    mariadb_port=3306 \
    mariadb_user=root \
    mariadb_pas=123456 \
    redis_host=trojan-panel-redis \
    redis_port=6379 \
    redis_pass=123456
COPY build/trojan-panel-linux-${get_arch} trojan-panel
# 国内环境开启以下注释 设置apk国内镜像
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add bash tzdata ca-certificates && \
    rm -rf /var/cache/apk/*
ENTRYPOINT ./trojan-panel \
    -host=\${mariadb_ip} \
    -port=\${mariadb_port} \
    -user=\${mariadb_user} \
    -password=\${mariadb_pas} \
    -redisHost=\${redis_host} \
    -redisPort=\${redis_port} \
    -redisPassword=\${redis_pass}
EOF

    docker buildx build --platform linux/"${get_arch}" -t jonssonyan/trojan-panel-linux-"${get_arch}" .
    if [[ "$?" == "0" ]]; then
      echo_content green "trojan-panel-linux-${get_arch}构建成功"
      echo_content skyBlue "开始推送trojan-panel-linux-${get_arch}"
      docker image tag jonssonyan/trojan-panel-linux-"${get_arch}":latest jonssonyan/trojan-panel:latest && \
      docker image push jonssonyan/trojan-panel:latest && \
      docker rmi -f jonssonyan/trojan-panel:latest
      if [[ "$?" == "0" ]]; then
        echo_content green "镜像名称：jonssonyan/trojan-panel:latest 架构：${get_arch}推送成功"
      else
        echo_content red "镜像名称：jonssonyan/trojan-panel:latest 架构：${get_arch}推送失败"
      fi

      if [[ ${trojan_panel_version} != "latest" ]]; then
        docker image tag jonssonyan/trojan-panel-linux-"${get_arch}":latest jonssonyan/trojan-panel:${trojan_panel_version} && \
        docker image push jonssonyan/trojan-panel:${trojan_panel_version} && \
        docker rmi -f jonssonyan/trojan-panel:${trojan_panel_version}
        if [[ "$?" == "0" ]]; then
          echo_content green "镜像名称：jonssonyan/trojan-panel:${trojan_panel_version} 架构：${get_arch}推送成功"
        else
          echo_content green "镜像名称：jonssonyan/trojan-panel:${trojan_panel_version} 架构：${get_arch}推送成功"
        fi
      fi
    else
      echo_content red "trojan-panel-linux-${get_arch}构建失败"
    fi
  done
}

init_var
main
