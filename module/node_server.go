package module

import "time"

type NodeServer struct {
	Id         *uint      `ddb:"id"`
	Name       *string    `ddb:"name"`
	Ip         *string    `ddb:"ip"`
	CreateTime *time.Time `ddb:"create_time"`
	UpdateTime *time.Time `ddb:"update_time"`
}
