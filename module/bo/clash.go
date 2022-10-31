package bo

type ClashConfigInterface interface {
	Vless | Vmess | Trojan | TrojanGo | Hysteria
}

type Vless struct {
}

type Vmess struct {
}

type Trojan struct {
}

type TrojanGo struct {
}

type Hysteria struct {
}

type ProxyGroups struct {
	Name      string   `yaml:"name"`
	ProxyType string   `yaml:"type"`
	Proxies   []string `yaml:"proxies"`
}

type ClashConfig struct {
	Proxies     []interface{} `yaml:"proxies"`
	ProxyGroups ProxyGroups   `yaml:"proxy-groups"`
}
