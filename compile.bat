::Windows amd64
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -ldflags="-H windowsgui -s -w" -o build/trojan-panel-win.exe
::Mac amd64
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -ldflags "-s -w" -o build/trojan-panel-mac
::Linux amd64
go install mvdan.cc/garble@latest
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
garble -literals build -o build/trojan-panel-linux-amd64
::Linux arm64
go install mvdan.cc/garble@latest
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=arm64
garble -literals build -o build/trojan-panel-linux-arm64