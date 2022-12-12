package vo

type PanelGroupVo struct {
	Quota        int     `json:"quota"`
	ResidualFlow int     `json:"residualFlow"`
	NodeCount    int     `json:"nodeCount"`
	ExpireTime   uint    `json:"expireTime"`
	AccountCount int     `json:"accountCount"`
	CpuUsed      float64 `json:"cpuUsed"`
	MemUsed      float64 `json:"memUsed"`
	DiskUsed     float64 `json:"diskUsed"`
}
