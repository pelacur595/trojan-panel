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
	Username    string `json:"username"`
	TrafficUsed string `json:"trafficUsed"`
}
