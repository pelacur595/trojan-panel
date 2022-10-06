package util

import (
	"errors"
	"github.com/go-ping/ping"
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
