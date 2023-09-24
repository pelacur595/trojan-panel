package model

import "time"

type NodeXray struct {
	Id             *uint      `ddb:"id"`
	Protocol       *string    `ddb:"protocol"`
	XrayFlow       *string    `ddb:"xray_flow"`
	XraySSMethod   *string    `ddb:"xray_ss_method"`
	RealityPbk     *string    `ddb:"reality_pbk"`
	Settings       *string    `ddb:"settings"`
	StreamSettings *string    `ddb:"stream_settings"`
	Tag            *string    `ddb:"tag"`
	Sniffing       *string    `ddb:"sniffing"`
	Allocate       *string    `ddb:"allocate"`
	CreateTime     *time.Time `ddb:"create_time"`
	UpdateTime     *time.Time `ddb:"update_time"`
}
