package util

import (
	"github.com/gin-gonic/gin"
	"strings"
	"trojan-panel/module/vo"
)

// GetCurrentAccount 获取当前用户
func GetCurrentAccount(c *gin.Context) *vo.AccountVo {
	// 解析token获取当前用户信息
	myClaims, err := ParseToken(GetToken(c))
	if err != nil {
		vo.Fail(err.Error(), c)
		return nil
	}
	accountVo := myClaims.AccountVo
	return &accountVo
}

// GetToken 获取token
func GetToken(c *gin.Context) string {
	tokenStr := c.Request.Header.Get("Authorization")
	if tokenStr == "" {
		return ""
	}
	return strings.SplitN(tokenStr, " ", 2)[1]
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

func IsAdmin(roleNames []string) bool {
	for _, item := range roleNames {
		if item == "admin" {
			return true
		}
	}
	return false
}
