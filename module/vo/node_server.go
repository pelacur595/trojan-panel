package vo

import "time"

type NodeServerVo struct {
	Id         uint      `json:"id"`
	Name       string    `json:"name"`
	Ip         string    `json:"ip"`
	GrpcPort   uint      `json:"grpcPort"`
	CreateTime time.Time `json:"createTime"`

	Status int `json:"status"`
}

type NodeServerPageVo struct {
	NodeServers []NodeServerVo `json:"nodeServers"`
	BaseVoPage
}

type NodeServerOneVo struct {
	Id         uint      `json:"id"`
	Name       string    `json:"name"`
	Ip         string    `json:"ip"`
	GrpcPort   uint      `json:"grpcPort"`
	CreateTime time.Time `json:"createTime"`
}

type NodeServerListVo struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type NodeServerGroupVo struct {
	CpuUsed  float32 `json:"cpuUsed"`
	MemUsed  float32 `json:"memUsed"`
	DiskUsed float32 `json:"diskUsed"`
}
