package vo

import "time"

type FileTaskVo struct {
	Id              uint      `json:"id"`
	Name            string    `json:"name"`
	Type            uint      `json:"type"`
	Status          int       `json:"status"`
	AccountUsername string    `json:"accountUsername"`
	CreateTime      time.Time `json:"create_time"`
}

type FileTaskPageVo struct {
	FileTaskVos []FileTaskVo `json:"fileTasks"`
	BaseVoPage
}
