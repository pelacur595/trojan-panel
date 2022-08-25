package module

import "time"

type Node struct {
	Id   *uint   `ddb:"id"`
	Name *string `ddb:"name"`
	Ip   *string `ddb:"ip"`
	Port *uint   `ddb:"port"`

	CreateTime *time.Time `ddb:"create_time"`
	UpdateTime *time.Time `ddb:"update_time"`
}
