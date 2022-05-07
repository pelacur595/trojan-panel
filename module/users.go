package module

import "time"

// 账户
type Users struct {
	Id                 *uint      `ddb:"id"`
	Password           *string    `ddb:"password"`
	Quota              *int       `ddb:"quota"`
	Download           *uint      `ddb:"download"`
	Upload             *uint      `ddb:"upload"`
	Username           *string    `ddb:"username"`
	Pass               *string    `ddb:"pass"`
	RoleId             *uint      `ddb:"role_id"`
	Deleted            *uint      `ddb:"deleted"`
	Email              *string    `ddb:"email"`
	ExpireTime         *uint      `ddb:"expire_time"`
	IpLimit            *uint      `ddb:"ip_limit"`
	UploadSpeedLimit   *uint      `ddb:"upload_speed_limit"`
	DownloadSpeedLimit *uint      `ddb:"download_speed_limit"`
	CreateTime         *time.Time `ddb:"create_time"`
	UpdateTime         *time.Time `ddb:"update_time"`
}
