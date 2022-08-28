package module

import "time"

type NodeTrojanGo struct {
	Id              *uint      `ddb:"id"`
	Sni             *string    `ddb:"sni"`
	MuxEnable       *uint      `ddb:"mux_enable"`
	WebsocketEnable *uint      `ddb:"websocket_enable"`
	WebsocketPath   *string    `ddb:"websocket_path"`
	SsEnable        *uint      `ddb:"ss_enable"`
	SsMethod        *string    `ddb:"ss_method"`
	SsPassword      *string    `ddb:"ss_password"`
	CreateTime      *time.Time `ddb:"create_time"`
	UpdateTime      *time.Time `ddb:"update_time"`
}
