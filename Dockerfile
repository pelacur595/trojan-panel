FROM golang:1.17

WORKDIR /tpdata/trojan-panel

COPY . /tpdata/trojan-panel/

RUN go build -o build/trojan-panel

ENTRYPOINT ["build/trojan-panel"]