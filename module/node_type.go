package module

import "time"

type NodeType struct {
	Id         *uint      `ddb:"id"`
	Name       *string    `ddb:"name"`
	Prefix     *string    `ddb:"prefix"`
	CreateTime *time.Time `ddb:"create_time"`
	UpdateTime *time.Time `ddb:"update_time"`
}
