package bo

type ClashConfigInterface interface {
	Vless | Vmess | Trojan | TrojanGo | Hysteria
}

type Vless struct {
	Name           string `yaml:"name"`
	Type           string `yaml:"type"`
	Server         string `yaml:"server"`
	Port           uint   `yaml:"port"`
	Uuid           string `yaml:"uuid"`
	Network        string `yaml:"network"`
	Tls            bool   `yaml:"tls"`
	Udp            bool   `yaml:"udp"`
	Flow           string `yaml:"flow"`
	SkipCertVerify bool   `yaml:"skip-cert-verify"`

	ServerName  string      `yaml:"servername"`
	RealityOpts RealityOpts `yaml:"reality-opts"`

	ClientFingerprint string `yaml:"client-fingerprint"`

	WsOpts WsOpts `yaml:"ws-opts"`
}

type RealityOpts struct {
	PublicKey string `yaml:"public-key"`
	ShortId   string `yaml:"short-id"`
}

type Vmess struct {
	Name              string `yaml:"name"`
	Type              string `yaml:"type"`
	Server            string `yaml:"server"`
	Port              uint   `yaml:"port"`
	Uuid              string `yaml:"uuid"`
	AlterId           uint   `yaml:"alterId"`
	Cipher            string `yaml:"cipher"`
	Tls               bool   `yaml:"tls"`
	Udp               bool   `yaml:"udp"`
	ClientFingerprint string `yaml:"client-fingerprint"`
	SkipCertVerify    bool   `yaml:"skip-cert-verify"`
	ServerName        string `yaml:"servername"`
	Network           string `yaml:"network"`
	WsOpts            WsOpts `yaml:"ws-opts"`
}

type Trojan struct {
	Name     string `yaml:"name"`
	Server   string `yaml:"server"`
	Port     uint   `yaml:"port"`
	Type     string `yaml:"type"`
	Password string `yaml:"password"`
	Udp      bool   `yaml:"udp"`
	Network  string `yaml:"network"`
	WsOpts   WsOpts `yaml:"ws-opts"`
}

type Shadowsocks struct {
	Name     string `yaml:"name"`
	Server   string `yaml:"server"`
	Port     uint   `yaml:"port"`
	Type     string `yaml:"type"`
	Cipher   string `yaml:"cipher"`
	Password string `yaml:"password"`
}

type TrojanGo struct {
	Name     string `yaml:"name"`
	Server   string `yaml:"server"`
	Port     uint   `yaml:"port"`
	Type     string `yaml:"type"`
	Password string `yaml:"password"`
	SNI      string `yaml:"sni"`
	Udp      bool   `yaml:"udp"`
	Network  string `yaml:"network"`
	WsOpts   WsOpts `yaml:"ws-opts"`
}

type Hysteria struct {
	Name     string `yaml:"name"`
	Server   string `yaml:"server"`
	Port     uint   `yaml:"port"`
	Type     string `yaml:"type"`
	AuthStr  string `yaml:"auth_str"`
	Obfs     string `yaml:"obfs"`
	Protocol string `yaml:"protocol"`
	Up       int    `yaml:"up"`
	Down     int    `yaml:"down"`
}

type WsOpts struct {
	Path    string        `yaml:"path"`
	Headers WsOptsHeaders `yaml:"headers"`
}

type WsOptsHeaders struct {
	Host string `yaml:"Host"`
}

type ProxyGroup struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Proxies []string `yaml:"proxies"`
}

type ClashConfig struct {
	Proxies     []interface{} `yaml:"proxies"`
	ProxyGroups []ProxyGroup  `yaml:"proxy-groups"`
}
