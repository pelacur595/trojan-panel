package router

import (
	"github.com/gin-gonic/gin"
	"trojan/api"
	"trojan/middleware"
)

func Router(router *gin.Engine) {
	router.Use(middleware.RateLimiterHandler(), middleware.LogHandler())
	trojanAuth := router.Group("/api/auth")
	{
		// 登录
		trojanAuth.POST("/login", api.Login)
		// 创建账户
		trojanAuth.POST("/register", api.Register)
		// 系统默认设置
		trojanAuth.GET("/setting", api.Setting)
	}
	router.Use(middleware.JWTHandler(), middleware.CasbinHandler())
	trojan := router.Group("/api")
	{
		dashboard := trojan.Group("/dashboard")
		{
			// 仪表板
			dashboard.GET("/panelGroup", api.PanelGroup)
		}
		user := trojan.Group("/users")
		{
			// 注销
			user.GET("/logout", api.LoginOut)
			// 查询单个账户
			user.GET("/selectUserById", api.SelectUserById)
			// 创建账户
			user.POST("/createUser", api.CreateUser)
			// 获取当前用户信息
			user.GET("/getUserInfo", api.GetUserInfo)
			// 分页查询账户
			user.GET("/selectUserPage", api.SelectUserPage)
			// 通过id删除账户
			user.POST("/deleteUserById", api.DeleteUserById)
			// 修改密码
			user.POST("/updateUserPassByUsername", api.UpdateUserPassByUsername)
			// 修改账户
			user.POST("/updateUserById", api.UpdateUserById)
		}
		role := trojan.Group("/role")
		{
			// 查询角色列表
			role.GET("/selectRoleList", api.SelectRoleList)
		}
		node := trojan.Group("/node")
		{
			// 根据id查询节点
			node.GET("/selectNodeById", api.SelectNodeById)
			// 创建节点
			node.POST("/createNode", api.CreateNode)
			// 分页查询节点
			node.GET("/selectNodePage", api.SelectNodePage)
			// 删除节点
			node.POST("/deleteNodeById", api.DeleteNodeById)
			// 更新节点
			node.POST("/updateNodeById", api.UpdateNodeById)
			// 获取节点二维码
			node.GET("/nodeQRCode", api.NodeQRCode)
			// 复制URL
			node.GET("/nodeURL", api.NodeURL)
		}
		nodeType := trojan.Group("/nodeType")
		{
			// 查询节点类型列表
			nodeType.GET("/selectNodeTypeList", api.SelectNodeTypeList)
		}
		system := trojan.Group("/system")
		{
			// 查询系统设置
			system.GET("/selectSystemByName", api.SelectSystemByName)
			// 更新系统配置
			system.POST("/updateSystemById", api.UpdateSystemById)
			// 上传静态网站文件
			system.POST("/uploadWebFile", api.UploadWebFile)
		}
		trojanGFW := trojan.Group("/trojan-gfw")
		{
			// 查看状态
			trojanGFW.GET("/status", api.TrojanGFWStatus)
			// 重启
			trojanGFW.POST("/restart", api.TrojanGFWRestart)
			// 停止
			trojanGFW.POST("/stop", api.TrojanGFWStop)
		}
		trojanGO := trojan.Group("/trojan-go")
		{
			// 查看状态
			trojanGO.GET("/status", api.TrojanGOStatus)
			// 重启
			trojanGO.POST("/restart", api.TrojanGORestart)
			// 停止
			trojanGO.POST("/stop", api.TrojanGOStop)
		}
	}
}
