package util

import (
	"github.com/go-ping/ping"
	"runtime"
)

func Ping(ip string) (int, error) {
	pingEr, err := ping.NewPinger(ip)
	if err != nil {
		panic(err)
	}
	pingEr.Count = 1
	if runtime.GOOS == "windows" {
		pingEr.SetPrivileged(true)
	}
	err = pingEr.Run()
	if err != nil {
		panic(err)
	}
	milliseconds := pingEr.Statistics().AvgRtt.Milliseconds()
	return int(milliseconds), nil
}
