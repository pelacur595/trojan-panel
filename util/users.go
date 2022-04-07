package util

import (
	"github.com/gin-gonic/gin"
	"strings"
	"trojan/module/constant"
	"trojan/module/vo"
)

func GetCurrentUser(c *gin.Context) *vo.UsersVo {
	// 解析token获取当前用户信息
	tokenStr := c.Request.Header.Get("Authorization")
	token := strings.SplitN(tokenStr, " ", 2)
	myClaims, err := ParseToken(token[1])
	if err != nil {
		vo.Fail(err.Error(), c)
		return nil
	}
	userVo := myClaims.UserVo
	return &userVo
}

func ToMB(b int) int {
	if b >= 0 {
		return b / 1024 / 1024
	} else {
		return -1
	}
}

func ToByte(b int) int {
	if b >= 0 {
		return b * 1024 * 1024
	} else {
		return -1
	}
}

func IsAdmin(id uint) bool {
	return constant.SYSADMIN == id || constant.ADMIN == id
}
