package module

import "time"

// 账户
type Users struct {
	Id         *uint      `ddb:"id"`
	Password   *string    `ddb:"password"`
	Quota      *int       `ddb:"quota"`
	Download   *uint      `ddb:"download"`
	Upload     *uint      `ddb:"upload"`
	Username   *string    `ddb:"username"`
	Pass       *string    `ddb:"pass"`
	RoleId     *uint      `ddb:"role_id"`
	Deleted    *uint      `ddb:"deleted"`
	ExpireTime *uint      `ddb:"expire_time"`
	CreateTime *time.Time `ddb:"create_time"`
	UpdateTime *time.Time `ddb:"update_time"`
}
