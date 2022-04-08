package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
	"trojan/dao"
	"trojan/module/constant"
	"trojan/module/vo"
	"trojan/util"
)

// casbin鉴权中间件
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		casbin, err := dao.Casbin()
		if err != nil {
			vo.Fail(err.Error(), c)
			c.Abort()
			return
		}
		// 检查账户权限
		path := c.Request.URL.Path
		split := strings.Split(path, "?")
		// 获取当前用户
		userVo := util.GetCurrentUser(c)
		roleVo, err := dao.SelectRoleById(&userVo.RoleId)
		if err != nil {
			vo.Fail(err.Error(), c)
			c.Abort()
			return
		}
		has, err := casbin.Enforce(roleVo.Name, split[0], c.Request.Method)
		if err != nil {
			vo.Fail(err.Error(), c)
			c.Abort()
			return
		}
		if !has {
			vo.Fail(constant.ForbiddenError, c)
			c.Abort()
			return
		}
		c.Next()
	}
}
