package vo

import (
	"time"
)

type UsersVo struct {
	Id         uint      `json:"id"`
	Quota      int       `json:"quota"`
	Download   uint      `json:"download"`
	Upload     uint      `json:"upload"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	RoleId     uint      `json:"roleId"`
	Deleted    uint      `json:"deleted"`
	ExpireTime uint      `json:"expireTime"`
	CreateTime time.Time `json:"createTime"`
}

type UsersPageVo struct {
	BaseVoPage
	Users []UsersVo `json:"users"`
}

type UsersLoginVo struct {
	Token string `json:"token"`
}

type UserInfo struct {
	Id        uint       `json:"id"`
	Username  string     `json:"username"`
	RoleNames []string   `json:"roleNames"`
	MenuList  []TreeNode `json:"menuList"`
}
