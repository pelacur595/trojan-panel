package module

import "time"

type System struct {
	Id            *uint      `ddb:"id"`
	Name          *string    `ddb:"name"`
	AccountConfig *string    `ddb:"account_config"` // 新用户设置
	EmailConfig   *string    `ddb:"email_config"`   // 系统邮箱设置
	CreateTime    *time.Time `ddb:"create_time"`
	UpdateTime    *time.Time `ddb:"update_time"`
}
