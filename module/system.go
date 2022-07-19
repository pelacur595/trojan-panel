package module

import "time"

type System struct {
	Id                 *uint      `ddb:"id"`
	Name               *string    `ddb:"name"`
	OpenRegister       *uint      `ddb:"open_register"`        // 开放注册 0/否 1/是
	RegisterQuota      *int       `ddb:"register_quota"`       // 注册用户默认配额 单位/MB
	RegisterExpireDays *uint      `ddb:"register_expire_days"` // 注册用户过期天数 单位/天
	ExpireWarnEnable   *uint      `ddb:"expire_warn_enable"`   // 是否开启到期警告 0/否 1/是
	ExpireWarnDay      *uint      `ddb:"expire_warn_day"`      // 到期警告 单位/天
	EmailEnable        *uint      `ddb:"email_enable"`         // 是否开启邮箱功能 0/否 1/是
	EmailHost          *string    `ddb:"email_host"`           // 系统邮箱设置-host
	EmailPort          *uint      `ddb:"email_port"`           // 系统邮箱设置-port
	EmailUsername      *string    `ddb:"email_username"`       // 系统邮箱设置-username
	EmailPassword      *string    `ddb:"email_password"`       // 系统邮箱设置-password
	CreateTime         *time.Time `ddb:"create_time"`
	UpdateTime         *time.Time `ddb:"update_time"`
}
