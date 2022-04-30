package middleware

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"trojan/module/constant"
	"trojan/module/vo"
)

var limit *limiter.Limiter

// 限流中间件
func RateLimiterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(limit, c.Writer, c.Request)
		if httpError != nil {
			vo.Fail(constant.RateLimiterError, c)
			c.Abort()
			return
		}
		c.Next()
	}
}

// 限流初始化
func init() {
	limit = tollbooth.NewLimiter(5, nil)
}
