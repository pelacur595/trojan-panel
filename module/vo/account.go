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
	Username   string `json:"username"`
	Pass       string `json:"pass"`
	Hash       string `json:"hash"`
	RoleId     string `json:"role_id"`
	Email      string `json:"email"`
	ExpireTime string `json:"expire_time"`
	Deleted    string `json:"deleted"`
	Quota      string `json:"quota"`
	Download   string `json:"download"`
	Upload     string `json:"upload"`
	CreateTime string `json:"create_time"`
}
