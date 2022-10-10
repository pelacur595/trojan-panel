package util

import (
	"errors"
	"fmt"
	"github.com/go-ping/ping"
	"github.com/sirupsen/logrus"
	"net"
	"runtime"
	"trojan-panel/module/constant"
)

func Ping(ip string) (int, error) {
	pingEr, err := ping.NewPinger(ip)
	if err != nil {
		return -1, errors.New(constant.SysError)
	}
	pingEr.Count = 3
	pingEr.Timeout = 3
	if runtime.GOOS == "windows" {
		pingEr.SetPrivileged(true)
	}
	err = pingEr.Run()
	if err != nil {
		return -1, errors.New(constant.SysError)
	}
	milliseconds := pingEr.Statistics().AvgRtt.Milliseconds()
	return int(milliseconds), nil
}

// IsPortAvailable 判断端口是否可用
func IsPortAvailable(port uint, network string) bool {
	address := fmt.Sprintf("127.0.0.1:%d", port)
	listener, err := net.Listen(network, address)
	if err != nil {
		logrus.Errorf("port %s is taken: %s \n", address, err)
		return false
	}

	defer listener.Close()
	return true
}
