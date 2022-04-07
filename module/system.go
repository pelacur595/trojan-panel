package module

import "time"

type System struct {
	Id                 *uint      `ddb:"id"`
	Name               *string    `ddb:"name"`
	OpenRegister       *uint      `ddb:"open_register"`        // 开放注册 0/否 1/是
	RegisterQuota      *int       `ddb:"register_quota"`       // 注册用户默认配额 单位/byte
	RegisterExpireDays *uint      `ddb:"register_expire_days"` // 注册用户过期天数 单位/天
	CreateTime         *time.Time `ddb:"create_time"`
	UpdateTime         *time.Time `ddb:"update_time"`
}
