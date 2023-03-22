package module

import "time"

type FileTask struct {
	Id         *uint      `ddb:"id"`
	Name       *string    `ddb:"name"`
	Path       *string    `ddb:"path"`
	Type       *uint      `ddb:"type"`
	Status     *int       `ddb:"status"`
	CreateTime *time.Time `ddb:"create_time"`
	UpdateTime *time.Time `ddb:"update_time"`
}
