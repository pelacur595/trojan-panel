# trojan-panel

Trojan Panel后端

# 编译命令

```
# Windows
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -ldflags="-H windowsgui -s -w" -o build/trojan-panel_win.exe

# Mac
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -ldflags "-s -w" -o build/trojan-panel_mac

# Linux
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-s -w" -o build/trojan-panel
# 加密编译（推荐）
go install mvdan.cc/garble@latest
garble -literals build -o build/trojan-panel
```

# Telegram讨论组

[Trojan Panel交流群](https://t.me/TrojanPanelGroup)

# 致谢

- [trojan](https://trojan-gfw.github.io/trojan/authenticator)
