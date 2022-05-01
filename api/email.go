package api

import (
	"github.com/gin-gonic/gin"
	"trojan/module/constant"
	"trojan/module/dto"
	"trojan/module/vo"
	"trojan/service"
)

func SendEmail(c *gin.Context) {
	var sendEmailDto dto.SendEmailDto
	_ = c.ShouldBindJSON(&sendEmailDto)
	if err := service.SendEmail(&sendEmailDto); err != nil {
		vo.Fail(constant.SendEmailError, c)
		return
	}
	vo.Success(nil, c)
}
