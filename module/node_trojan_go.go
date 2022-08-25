package module

type NodeTrojanGo struct {
	Sni             *string `ddb:"sni"`
	Type            *uint   `ddb:"type"`
	WebsocketEnable *uint   `ddb:"websocket_enable"`
	WebsocketPath   *string `ddb:"websocket_path"`
	SsEnable        *uint   `ddb:"ss_enable"`
	SsMethod        *string `ddb:"ss_method"`
	SsPassword      *string `ddb:"ss_password"`
}
