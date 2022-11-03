package bo

type ClashConfigInterface interface {
	Vless | Vmess | Trojan | TrojanGo | Hysteria
}

type Vless struct {
}

type Vmess struct {
	Name      string `yaml:"name"`
	Server    string `yaml:"server"`
	Port      uint   `yaml:"port"`
	VmessType string `yaml:"type"`
	Uuid      string `yaml:"uuid"`
	AlterId   uint   `yaml:"alterId"`
	Cipher    string `yaml:"cipher"`
	Udp       bool   `yaml:"udp"`
	Tls       bool   `yaml:"tls"`
	Network   string `yaml:"network"`
	WsOpts    WsOpts `yaml:"ws-opts"`
}

type Trojan struct {
	Name       string `yaml:"name"`
	Server     string `yaml:"server"`
	Port       uint   `yaml:"port"`
	TrojanType string `yaml:"type"`
	Password   string `yaml:"password"`
	SNI        string `yaml:"sni"`
	Udp        bool   `yaml:"udp"`
	Network    string `yaml:"network"`
	WsOpts     WsOpts `yaml:"ws-opts"`
}

type TrojanGo struct {
	Name       string `yaml:"name"`
	Server     string `yaml:"server"`
	Port       uint   `yaml:"port"`
	TrojanType string `yaml:"type"`
	Password   string `yaml:"password"`
	SNI        string `yaml:"sni"`
	Udp        bool   `yaml:"udp"`
	Network    string `yaml:"network"`
	WsOpts     WsOpts `yaml:"ws-opts"`
}

type Hysteria struct {
}

type WsOpts struct {
	Path          string        `yaml:"path"`
	WsOptsHeaders WsOptsHeaders `yaml:"headers"`
}

type WsOptsHeaders struct {
	Host string `yaml:"Host"`
}

type ProxyGroup struct {
	Name      string   `yaml:"name"`
	ProxyType string   `yaml:"type"`
	Proxies   []string `yaml:"proxies"`
}

type ClashConfig struct {
	Proxies     []interface{} `yaml:"proxies"`
	ProxyGroups []ProxyGroup  `yaml:"proxy-groups"`
}
