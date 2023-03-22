package vo

import "time"

type FileTaskVo struct {
	Id         uint      `ddb:"id"`
	Name       string    `ddb:"name"`
	Type       uint      `ddb:"type"`
	Status     int       `ddb:"status"`
	CreateTime time.Time `ddb:"create_time"`
}

type FileTaskPageVo struct {
	FileTaskVos []FileTaskVo `json:"fileTasks"`
	BaseVoPage
}
