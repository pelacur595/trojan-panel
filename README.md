# trojan-panel

Trojan Panel后端

# 编译命令

```
# Windows
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -ldflags="-H windowsgui -s -w" -o build/trojan-panel-win.exe

# Mac
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -ldflags "-s -w" -o build/trojan-panel-mac

# Linux amd64
go install mvdan.cc/garble@latest
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
# 普通编译
go build -ldflags "-s -w" -o build/trojan-panel-linux-amd64
# 加密编译（推荐）
garble -literals build -o build/trojan-panel-linux-amd64
```

# Telegram讨论组

[Trojan Panel交流群](https://t.me/TrojanPanelGroup)

# 致谢

- [trojan-gfw](https://github.com/trojan-gfw/trojan)
- [trojan-go](https://github.com/p4gefau1t/trojan-go)
