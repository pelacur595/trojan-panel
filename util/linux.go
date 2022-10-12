package util

import (
	"errors"
	"github.com/go-ping/ping"
	"github.com/sirupsen/logrus"
	"net"
	"runtime"
	"time"
	"trojan-panel/module/constant"
)

func Ping(ip string) (int, error) {
	pingEr, err := ping.NewPinger(ip)
	if err != nil {
		return -1, errors.New(constant.SysError)
	}
	pingEr.Count = 1
	pingEr.Timeout = 2 * time.Second
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
	if network == "tcp" {
		listener, err := net.ListenTCP(network, &net.TCPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: int(port),
		})
		defer listener.Close()
		if err != nil {
			logrus.Errorf("port %d is taken err: %s\n", port, err)
			return false
		}
	}
	if network == "udp" {
		listener, err := net.ListenUDP("udp", &net.UDPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: int(port),
		})
		defer listener.Close()
		if err != nil {
			logrus.Errorf("port %d is taken err: %s\n", port, err)
			return false
		}
	}
	return true
}
