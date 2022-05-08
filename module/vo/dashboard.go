package vo

type PanelGroupVo struct {
	Quota        int  `json:"quota"`
	ResidualFlow int  `json:"residualFlow"`
	NodeCount    int  `json:"nodeCount"`
	ExpireTime   uint `json:"expireTime"`
	UserCount    int  `json:"userCount"`
	OnLine       int  `json:"onLine"`
}
