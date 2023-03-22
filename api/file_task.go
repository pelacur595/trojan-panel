package api

import (
	"github.com/gin-gonic/gin"
	"trojan-panel/module/constant"
	"trojan-panel/module/dto"
	"trojan-panel/module/vo"
	"trojan-panel/service"
)

func SelectFileTaskPage(c *gin.Context) {
	var fileTaskPageDto dto.FileTaskPageDto
	_ = c.ShouldBindQuery(&fileTaskPageDto)
	if err := validate.Struct(&fileTaskPageDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	page, err := service.SelectFileTaskPage(fileTaskPageDto.Type, fileTaskPageDto.PageNum, fileTaskPageDto.PageSize)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(page, c)
}

func DeleteFileTaskById(c *gin.Context) {
	var accountRequiredIdDto dto.RequiredIdDto
	_ = c.ShouldBindJSON(&accountRequiredIdDto)
	if err := validate.Struct(&accountRequiredIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	if err := service.DeleteFileTaskById(accountRequiredIdDto.Id); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}
