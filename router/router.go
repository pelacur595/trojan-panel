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
		// Hysteria api
		trojanAuth.POST("/hysteria", api.HysteriaApi)
	}
	router.Use(middleware.JWTHandler(), middleware.CasbinHandler())
	trojan := router.Group("/api")
	{
		dashboard := trojan.Group("/dashboard")
		{
			// 仪表板
			dashboard.GET("/panelGroup", api.PanelGroup)
			// 流量排行榜
			dashboard.GET("/trafficRank", api.TrafficRank)
		}
		account := trojan.Group("/account")
		{
			// 注销
			account.POST("/logout", api.Logout)
			// 查询单个账户
			account.GET("/selectAccountById", api.SelectAccountById)
			// 创建账户
			account.POST("/createAccount", api.CreateAccount)
			// 获取当前用户信息
			account.GET("/getAccountInfo", api.GetAccountInfo)
			// 分页查询账户
			account.GET("/selectAccountPage", api.SelectAccountPage)
			// 通过id删除账户
			account.POST("/deleteAccountById", api.DeleteAccountById)
			// 修改个人信息
			account.POST("/updateAccountProfile", api.UpdateAccountProfile)
			// 修改账户
			account.POST("/updateAccountById", api.UpdateAccountById)
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
			node.POST("/nodeQRCode", api.NodeQRCode)
			// 复制URL
			node.POST("/nodeURL", api.NodeURL)
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
		blackList := trojan.Group("/blackList")
		{
			// 分页查询黑名单
			blackList.GET("/selectBlackListPage", api.SelectBlackListPage)
			// 删除黑名单
			blackList.POST("/deleteBlackListByIp", api.DeleteBlackListByIp)
			// 创建黑名单
			blackList.POST("/createBlackList", api.CreateBlackList)
		}
		emailRecord := trojan.Group("/emailRecord")
		{
			// 查询邮件发送记录
			emailRecord.GET("/selectEmailRecordPage", api.SelectEmailRecordPage)
		}
	}
}
