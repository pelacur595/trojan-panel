package module

import "time"

// Users 账户
type Users struct {
	Id                 *uint      `ddb:"id"`
	Username           *string    `ddb:"username"`
	Pass               *string    `ddb:"pass"`
	RoleId             *uint      `ddb:"role_id"`
	Email              *string    `ddb:"email"`
	ExpireTime         *uint      `ddb:"expire_time"`
	Deleted            *uint      `ddb:"deleted"`
	Quota              *int       `ddb:"quota"`
	Download           *int       `ddb:"download"`
	Upload             *int       `ddb:"upload"`
	IpLimit            *uint      `ddb:"ip_limit"`
	UploadSpeedLimit   *uint      `ddb:"upload_speed_limit"`
	DownloadSpeedLimit *uint      `ddb:"download_speed_limit"`
	CreateTime         *time.Time `ddb:"create_time"`
	UpdateTime         *time.Time `ddb:"update_time"`
}
