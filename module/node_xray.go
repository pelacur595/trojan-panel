package module

import "time"

type NodeXray struct {
	Id             *uint      `ddb:"id"`
	Protocol       *string    `ddb:"protocol"`
	Settings       *string    `ddb:"settings"`
	StreamSettings *string    `ddb:"stream_settings"`
	Tag            *string    `ddb:"tag"`
	Sniffing       *string    `ddb:"sniffing"`
	Allocate       *string    `ddb:"allocate"`
	CreateTime     *time.Time `ddb:"create_time"`
	UpdateTime     *time.Time `ddb:"update_time"`
}
