go install mvdan.cc/garble@latest
::Windows amd64
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -ldflags="-H windowsgui -s -w" -o build/trojan-panel-win-amd64.exe
::Mac amd64
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -ldflags "-s -w" -o build/trojan-panel-mac-amd64
::Linux amd64
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
garble -literals build -o build/trojan-panel-linux-amd64
::Linux arm64
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=arm64
garble -literals build -o build/trojan-panel-linux-arm64