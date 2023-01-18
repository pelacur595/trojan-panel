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

type XrayTemplate struct {
	Flow           string `json:"flow"`
	ConfigTemplate string `json:"configTemplate"`
}

type XrayConfigBo struct {
	Log       TypeMessage `json:"log"`
	API       TypeMessage `json:"api"`
	DNS       TypeMessage `json:"dns"`
	Routing   TypeMessage `json:"routing"`
	Policy    TypeMessage `json:"policy"`
	Inbounds  TypeMessage `json:"inbounds"`
	Outbounds TypeMessage `json:"outbounds"`
	Transport TypeMessage `json:"transport"`
	Stats     TypeMessage `json:"stats"`
	Reverse   TypeMessage `json:"reverse"`
	FakeDNS   TypeMessage `json:"fakeDns"`
}
