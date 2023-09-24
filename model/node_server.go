package model

import "time"

type NodeServer struct {
	Id         *uint      `ddb:"id" json:"id"`
	Name       *string    `ddb:"name" json:"name"`
	Ip         *string    `ddb:"ip" json:"ip"`
	GrpcPort   *uint      `ddb:"grpc_port" json:"grpcPort"`
	CreateTime *time.Time `ddb:"create_time" json:"createTime"`
	UpdateTime *time.Time `ddb:"update_time" json:"updateTime"`
}
