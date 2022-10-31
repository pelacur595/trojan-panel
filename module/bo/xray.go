package bo

type StreamSettings struct {
	Network    string     `json:"network"`
	Security   string     `json:"security"`
	WsSettings WsSettings `json:"wsSettings"`
}

type WsSettings struct {
	Path string `json:"path"`
	Host string `json:"host"`
}

type Settings struct {
	Encryption string `json:"encryption"`
}
