package module

import "time"

type Node struct {
	Id               *uint      `ddb:"id"`
	Name             *string    `ddb:"name"`
	Ip               *string    `ddb:"ip"`
	Port             *uint      `ddb:"port"`
	Type             *uint      `ddb:"type"`
	WebsocketEnable  *uint      `ddb:"websocket_enable"`
	WebsocketPath    *string    `ddb:"websocket_path"`
	SsEnable         *uint      `ddb:"ss_enable"`
	SsMethod         *string    `ddb:"ss_method"`
	SsPassword       *string    `ddb:"ss_password"`
	HysteriaProtocol *string    `ddb:"hysteria_protocol"`
	CreateTime       *time.Time `ddb:"create_time"`
	UpdateTime       *time.Time `ddb:"update_time"`
}
