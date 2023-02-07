package module

import "time"

type NodeServer struct {
	Id         *uint      `ddb:"id"`
	Name       *string    `ddb:"name"`
	Ip         *string    `ddb:"ip"`
	GrpcPort   *uint      `ddb:"grpc_port"`
	CreateTime *time.Time `ddb:"create_time"`
	UpdateTime *time.Time `ddb:"update_time"`
}
