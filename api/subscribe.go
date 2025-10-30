package api

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"github.com/gin-gonic/gin"
	"trojan-panel/model/constant"
	"trojan-panel/model/dto"
	"trojan-panel/model/vo"
	"trojan-panel/service"
)

// getSubscriptionURL 生成订阅链接（支持固定端口）
func getSubscriptionURL(c *gin.Context, token string) string {
	// 获取请求的 scheme (http/https)
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	
	// 获取域名（不包含端口）
	host := c.Request.Host
	// 移除端口号
	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}
	
	// 检查是否设置了固定的订阅端口
	subscriptionPort := os.Getenv("SUBSCRIPTION_PORT")
	if subscriptionPort != "" {
		// 使用固定端口
		return fmt.Sprintf("%s://%s:%s/api/auth/subscribe/%s", scheme, host, subscriptionPort, token)
	}
	
	// 如果没有设置固定端口，使用当前访问的端口（保持原有逻辑兼容性）
	return fmt.Sprintf("%s://%s/api/auth/subscribe/%s", scheme, c.Request.Host, token)
}

// ClashSubscribe 获取Clash订阅地址
func ClashSubscribe(c *gin.Context) {
	accountVo := service.GetCurrentAccount(c)
	password, err := service.SelectConnectPassword(&accountVo.Id, &accountVo.Username)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	
	token := base64.StdEncoding.EncodeToString([]byte(password))
	subscribeURL := getSubscriptionURL(c, token)
	
	vo.Success(subscribeURL, c)
}

// ClashSubscribeForSb 获取指定人的Clash订阅地址
func ClashSubscribeForSb(c *gin.Context) {
	var accountRequiredIdDto dto.RequiredIdDto
	_ = c.ShouldBindQuery(&accountRequiredIdDto)
	if err := validate.Struct(&accountRequiredIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	password, err := service.SelectConnectPassword(accountRequiredIdDto.Id, nil)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	
	token := base64.StdEncoding.EncodeToString([]byte(password))
	subscribeURL := getSubscriptionURL(c, token)
	
	vo.Success(subscribeURL, c)
}

// Subscribe 订阅
func Subscribe(c *gin.Context) {
	token := c.Param("token")
	//userAgent := c.Request.Header.Get("User-Agent")
	tokenDecode, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	pass := string(tokenDecode)
	//if strings.HasPrefix(userAgent, constant.ClashforWindows) {
	account, userInfo, clashConfigYaml, systemConfig, err := service.SubscribeClash(pass)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	result := fmt.Sprintf(`%s
%s`, string(clashConfigYaml), systemConfig.ClashRule)
	c.Header("content-disposition", fmt.Sprintf("attachment; filename=%s.yaml", *account.Username))
	c.Header("profile-update-interval", "12")
	c.Header("subscription-userinfo", userInfo)
	c.String(200, result)
	return
	//}
	//vo.Fail("This client is not supported", c)
}
