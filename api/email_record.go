package api

import (
	"github.com/gin-gonic/gin"
	"trojan/module"
	"trojan/module/constant"
	"trojan/module/dto"
	"trojan/module/vo"
	"trojan/service"
)

func SelectEmailRecordPage(c *gin.Context) {
	var emailRecordPageDto dto.EmailRecordPageDto
	_ = c.ShouldBindQuery(&emailRecordPageDto)
	if err := validate.Struct(&emailRecordPageDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	emailRecordPageVo, err := service.SelectEmailRecordPage(
		emailRecordPageDto.ToEmail, emailRecordPageDto.State,
		emailRecordPageDto.PageNum, emailRecordPageDto.PageSize)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(emailRecordPageVo, c)
}

func CreateEmailRecord(c *gin.Context) {
	var emailRecordCreateDto dto.EmailRecordCreateDto
	_ = c.ShouldBindQuery(&emailRecordCreateDto)
	if err := validate.Struct(&emailRecordCreateDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	emailRecord := module.EmailRecord{
		ToEmail: emailRecordCreateDto.ToEmail,
		Subject: emailRecordCreateDto.Subject,
		Content: emailRecordCreateDto.Content,
	}
	if err := service.CreateEmailRecord(&emailRecord); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func UpdateEmailRecordById(c *gin.Context) {
	var emailRecordUpdateDto dto.EmailRecordUpdateDto
	_ = c.ShouldBindQuery(&emailRecordUpdateDto)
	if err := validate.Struct(&emailRecordUpdateDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	emailRecord := module.EmailRecord{
		Id:    emailRecordUpdateDto.Id,
		State: emailRecordUpdateDto.State,
	}
	if err := service.UpdateEmailRecordById(&emailRecord); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}
