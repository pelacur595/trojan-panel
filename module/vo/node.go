package vo

import "time"

// NodeVo 查询分页Node对象
type NodeVo struct {
	Id           uint      `json:"id"`
	NodeServerId uint      `json:"nodeServerId"`
	NodeSubId    uint      `json:"nodeSubId"`
	NodeTypeId   uint      `json:"nodeTypeId"`
	Name         string    `json:"name"`
	Domain       string    `json:"domain"`
	Port         uint      `json:"port"`
	CreateTime   time.Time `json:"createTime"`

	Status int `json:"status"`
}

type NodePageVo struct {
	Nodes []NodeVo `json:"nodes"`
	BaseVoPage
}

// NodeOneVo 查询单个Node对象
type NodeOneVo struct {
	Id           uint      `json:"id"`
	NodeServerId uint      `json:"nodeServerId"`
	NodeSubId    uint      `json:"nodeSubId"`
	NodeTypeId   uint      `json:"nodeTypeId"`
	Name         string    `json:"name"`
	Domain       string    `json:"domain"`
	Port         uint      `json:"port"`
	CreateTime   time.Time `json:"createTime"`

	Password string `json:"password"`
	Uuid     string `json:"uuid"`
	AlterId  int    `json:"alterId"`

	XrayProtocol             string                   `json:"xrayProtocol"`
	XrayFlow                 string                   `json:"xrayFlow"`
	XraySSMethod             string                   `json:"xraySSMethod"`
	RealityPbk               string                   `json:"realityPbk"`
	XraySettings             string                   `json:"xraySettings"`
	XraySettingEntity        XraySettingEntity        `json:"xraySettingsEntity"`
	XrayStreamSettingsEntity XrayStreamSettingsEntity `json:"xrayStreamSettingsEntity"`
	XrayTag                  string                   `json:"xrayTag"`
	XraySniffing             string                   `json:"xraySniffing"`
	XrayAllocate             string                   `json:"xrayAllocate"`
	TrojanGoSni              string                   `json:"trojanGoSni"`
	TrojanGoMuxEnable        uint                     `json:"trojanGoMuxEnable"`
	TrojanGoWebsocketEnable  uint                     `json:"trojanGoWebsocketEnable"`
	TrojanGoWebsocketPath    string                   `json:"trojanGoWebsocketPath"`
	TrojanGoWebsocketHost    string                   `json:"trojanGoWebsocketHost"`
	TrojanGoSsEnable         uint                     `json:"trojanGoSsEnable"`
	TrojanGoSsMethod         string                   `json:"trojanGoSsMethod"`
	TrojanGoSsPassword       string                   `json:"trojanGoSsPassword"`
	HysteriaProtocol         string                   `json:"hysteriaProtocol"`
	HysteriaUpMbps           int                      `json:"hysteriaUpMbps"`
	HysteriaDownMbps         int                      `json:"hysteriaDownMbps"`
	NaiveProxyUsername       string                   `json:"naiveProxyUsername"`
}

type XrayStreamSettingsEntity struct {
	Network         string                                  `json:"network"`
	Security        string                                  `json:"security"`
	RealitySettings XrayStreamSettingsRealitySettingsEntity `json:"realitySettings"`
	WsSettings      XrayStreamSettingsWsSettingsEntity      `json:"wsSettings"`
}

type XraySettingEntity struct {
	Fallbacks []XrayFallback `json:"fallbacks"`
	Network   string         `json:"network"`
}

type XrayFallback struct {
	Name *string `json:"name"`
	Alpn *string `json:"alpn"`
	Path *string `json:"path"`
	Dest any     `json:"dest"`
	Xver *uint   `json:"xver"`
}

type XrayStreamSettingsRealitySettingsEntity struct {
	Dest        string   `json:"dest"`
	Xver        int      `json:"xver"`
	ServerNames []string `json:"serverNames"`
	PrivateKey  string   `json:"privateKey"`
	ShortIds    []string `json:"shortIds"`
}

type XrayStreamSettingsWsSettingsEntity struct {
	Path string `json:"path"`
}
