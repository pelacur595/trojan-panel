# trojan-panel

Trojan Panel后端

# 编译命令

```
# Windows
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o build/trojan-panel_dev.exe
go build -ldflags="-H windowsgui" -o build/trojan-panel.exe

# Mac
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -o build/trojan-panel_mac

# Linux
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o build/trojan-panel
```

# Telegram讨论组

[Trojan Panel交流群](https://t.me/TrojanPanelGroup)

# 致谢

- [trojan](https://trojan-gfw.github.io/trojan/authenticator)
