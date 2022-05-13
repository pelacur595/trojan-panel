package util

import "trojan/core"

// 在线人数
func OnLine(ip string) (uint, error) {
	api := core.TrojanGoApi()
	onLine, err := api.OnLine(ip)
	if err != nil {
		return 0, err
	}
	return onLine, nil
}
