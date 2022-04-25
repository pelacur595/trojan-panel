FROM golang:1.17

WORKDIR /tpdata/trojan-panel

COPY . .

RUN go build -ldflags "-s -w" -o build/trojan-panel

ENTRYPOINT ["./build/trojan-panel"]