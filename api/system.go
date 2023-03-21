package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"trojan-panel/module/constant"
	"trojan-panel/module/dto"
	"trojan-panel/module/vo"
	"trojan-panel/service"
	"trojan-panel/util"
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
		RegisterEnable:     systemVo.RegisterEnable,
		RegisterQuota:      systemVo.RegisterQuota,
		RegisterExpireDays: systemVo.RegisterExpireDays,
		TrafficRankEnable:  systemVo.TrafficRankEnable,
		EmailEnable:        systemVo.EmailEnable,
		SystemName:         systemVo.SystemName,
	}
	vo.Success(settingVo, c)
}

func UpdateSystemById(c *gin.Context) {
	var systemUpdateDto dto.SystemUpdateDto
	_ = c.ShouldBindJSON(&systemUpdateDto)
	if err := validate.Struct(&systemUpdateDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	if len(systemUpdateDto.FileList) > 0 {
		if err := util.DownloadFile(systemUpdateDto.FileList[0].Url, constant.LogoImagePath); err != nil {
			vo.Fail(constant.SysError, c)
			return
		}
	}
	if err := service.UpdateSystemById(systemUpdateDto); err != nil {
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
	if err := util.RemoveSubFile(constant.WebFilePath); err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	// 保存文件
	webFile := fmt.Sprintf("%s/%s", constant.WebFilePath, constant.WebFileName)
	if err := c.SaveUploadedFile(file, webFile); err != nil {
		vo.Fail(constant.FileUploadError, c)
		return
	}
	// 解压webfile
	if err := util.Unzip(webFile, constant.WebFilePath); err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	vo.Success(nil, c)
}

func GetLogo(c *gin.Context) {
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename=logo.png")
	c.File(constant.LogoImagePath)
}
