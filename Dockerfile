FROM golang:1.16

WORKDIR /tpdata/trojan-panel

COPY . /tpdata/trojan-panel/

RUN go build -o build/trojan-panel

ENTRYPOINT ["build/trojan-panel"]