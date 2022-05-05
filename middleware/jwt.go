package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
	"trojan/dao"
	"trojan/dao/redis"
	"trojan/module/constant"
	"trojan/module/vo"
	"trojan/util"
)

// jwt认证中间件
func JWTHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Token放在Header的Authorization中，并使用Bearer开头
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			vo.Fail(constant.UnauthorizedError, c)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			vo.Fail(constant.IllegalTokenError, c)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		myClaims, err := util.ParseToken(parts[1])
		if err != nil {
			vo.Fail(err.Error(), c)
			c.Abort()
			return
		}
		get := redis.Client.String.
			Get(fmt.Sprintf("trojan-panel:token:%s", myClaims.UserVo.Username))
		result, err := get.String()
		if err != nil || result == "" {
			vo.Fail(constant.IllegalTokenError, c)
			c.Abort()
			return
		}

		// IP黑名单
		ip := c.ClientIP()
		get = redis.Client.String.Get(ip)
		result, err = get.String()
		if err != nil {
			vo.Fail(constant.IllegalTokenError, c)
			c.Abort()
			return
		}
		if result != "" {
			redis.Client.String.Set(fmt.Sprintf("trojan-panel:black-list:%s", ip), "in-black-list", time.Hour.Milliseconds()/1000)
			vo.Fail(constant.BlackListError, c)
			c.Abort()
			return
		} else {
			ipCount, err := dao.CountBlackListByIp(&ip)
			if err != nil {
				vo.Fail(err.Error(), c)
				c.Abort()
				return
			}
			if ipCount > 0 {
				redis.Client.String.Set(fmt.Sprintf("trojan-panel:black-list:%s", ip), "in-black-list", time.Hour.Milliseconds()/1000)
				vo.Fail(constant.BlackListError, c)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
