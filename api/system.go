package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"trojan/module"
	"trojan/module/constant"
	"trojan/module/dto"
	"trojan/module/vo"
	"trojan/service"
	"trojan/util"
)

func SelectSystemByName(c *gin.Context) {
	name := constant.SystemName
	systemVo, err := service.SelectSystemByName(&name)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(systemVo, c)
}

func Setting(c *gin.Context) {
	name := constant.SystemName
	systemVo, err := service.SelectSystemByName(&name)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	settingVo := vo.SettingVo{
		OpenRegister: systemVo.OpenRegister,
	}
	vo.Success(settingVo, c)
}

func UpdateSystemById(c *gin.Context) {
	var systemUpdateDto dto.SystemUpdateDto
	_ = c.ShouldBindJSON(&systemUpdateDto)
	system := module.System{
		Id:                 systemUpdateDto.Id,
		OpenRegister:       systemUpdateDto.OpenRegister,
		RegisterQuota:      systemUpdateDto.RegisterQuota,
		RegisterExpireDays: systemUpdateDto.RegisterExpireDays,
	}
	if err := service.UpdateSystemById(&system); err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	vo.Success(nil, c)
}

func UploadWebFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	// 文件大小 10MB
	if file.Size > 1024*1024*10 {
		vo.Fail(constant.FileSizeTooBig, c)
		return
	}
	// 文件后缀.zip
	if !strings.HasSuffix(file.Filename, ".zip") {
		vo.Fail(constant.FileFormatError, c)
		return
	}
	// 删除webfile文件夹内的所有文件
	if err = util.RemoveSubFile(constant.WebFilePath); err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	// 保存文件
	webFile := fmt.Sprintf("%s/%s", constant.WebFilePath, constant.WebFileName)
	if err = c.SaveUploadedFile(file, webFile); err != nil {
		vo.Fail(constant.FileUploadError, c)
		return
	}
	// 解压webfile
	if err = util.Unzip(webFile, constant.WebFilePath); err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	vo.Success(nil, c)
}
