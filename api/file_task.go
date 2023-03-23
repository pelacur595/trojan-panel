package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
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
	page, err := service.SelectFileTaskPage(fileTaskPageDto.Type, fileTaskPageDto.AccountUsername, fileTaskPageDto.PageNum, fileTaskPageDto.PageSize)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(page, c)
}

func DeleteFileTaskById(c *gin.Context) {
	var requiredIdDto dto.RequiredIdDto
	_ = c.ShouldBindJSON(&requiredIdDto)
	if err := validate.Struct(&requiredIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	if err := service.DeleteFileTaskById(requiredIdDto.Id); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

// DownloadFileTask 下载文件任务的文件
func DownloadFileTask(c *gin.Context) {
	var requiredIdDto dto.RequiredIdDto
	_ = c.ShouldBindJSON(&requiredIdDto)
	if err := validate.Struct(&requiredIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	fileTask, err := service.SelectFileTaskById(requiredIdDto.Id)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}

	if fileTask == nil || *fileTask.Status != constant.TaskSuccess {
		vo.Fail(constant.FileTaskNotSuccess, c)
		return
	}

	// 打开文件
	file, err := os.Open(*fileTask.Path)
	if err != nil {
		vo.Fail(constant.FileNotExist, c)
		return
	}
	defer file.Close()

	// 设置响应头，以便在文件下载时使用该文件名
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", *fileTask.Name))
	// 设置响应类型为二进制流
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	// 传输文件到客户端
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	//vo.Success(nil, c)
}

func DownloadCsvTemplate(c *gin.Context) {
	var templateRequiredIdDto dto.RequiredIdDto
	_ = c.ShouldBindJSON(&templateRequiredIdDto)
	if err := validate.Struct(&templateRequiredIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	if *templateRequiredIdDto.Id == constant.TaskTypeAccountExport {
		c.Header("Content-Disposition", "attachment; filename=AccountTemplate.csv")
		c.File(constant.ExcelAccountTemplate)
		return
	} else if *templateRequiredIdDto.Id == constant.TaskTypeNodeServerExport {
		c.Header("Content-Disposition", "attachment; filename=NodeServerTemplate.csv")
		c.File(constant.ExcelNodeServerTemplate)
		return
	}
	vo.Fail(constant.FileNotExist, c)
}
