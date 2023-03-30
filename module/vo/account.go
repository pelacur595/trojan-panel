package vo

import (
	"time"
)

type AccountVo struct {
	Id         uint      `json:"id"`
	Quota      int       `json:"quota"`
	Download   int       `json:"download"`
	Upload     int       `json:"upload"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	RoleId     uint      `json:"roleId"`
	Deleted    uint      `json:"deleted"`
	ExpireTime uint      `json:"expireTime"`
	CreateTime time.Time `json:"createTime"`
	Roles      []string  `json:"roles"`
}

type AccountPageVo struct {
	BaseVoPage
	AccountVos []AccountVo `json:"accounts"`
}

type AccountLoginVo struct {
	Token string `json:"token"`
}

type AccountInfo struct {
	Id       uint     `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

type AccountTrafficRankVo struct {
	Username    string `json:"username" ddb:"username"`
	TrafficUsed string `json:"trafficUsed" ddb:"trafficUsed"`
}

type AccountExportVo struct {
	Username   string `json:"username" ddb:"username"`
	Pass       string `json:"pass" ddb:"pass"`
	Hash       string `json:"hash" ddb:"hash"`
	RoleId     string `json:"role_id" ddb:"role_id"`
	Email      string `json:"email" ddb:"email"`
	ExpireTime string `json:"expire_time" ddb:"expire_time"`
	Deleted    string `json:"deleted" ddb:"deleted"`
	Quota      string `json:"quota" ddb:"quota"`
	Download   string `json:"download" ddb:"download"`
	Upload     string `json:"upload" ddb:"upload"`
	CreateTime string `json:"create_time" ddb:"create_time"`
}

type CaptureVo struct {
	CaptchaId  string `json:"captchaId"`
	CaptchaImg string `json:"captchaImg"`
}
