package api

import (
	"github.com/gin-gonic/gin"
	"trojan-panel/module"
	"trojan-panel/module/constant"
	"trojan-panel/module/dto"
	"trojan-panel/module/vo"
	"trojan-panel/service"
)

func SelectNodeServerById(c *gin.Context) {
	var nodeServerRequireIdDto dto.RequiredIdDto
	_ = c.ShouldBindQuery(&nodeServerRequireIdDto)
	if err := validate.Struct(&nodeServerRequireIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	nodeServer, err := service.SelectNodeServerById(nodeServerRequireIdDto.Id)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	nodeServerOneVo := vo.NodeServerOneVo{
		Id:         *nodeServer.Id,
		Name:       *nodeServer.Name,
		Ip:         *nodeServer.Ip,
		CreateTime: *nodeServer.CreateTime,
	}
	vo.Success(nodeServerOneVo, c)
}

func CreateNodeServer(c *gin.Context) {
	var nodeServerCreateDto dto.NodeServerCreateDto
	_ = c.ShouldBindJSON(&nodeServerCreateDto)
	if err := validate.Struct(&nodeServerCreateDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	nodeServer := module.NodeServer{Name: nodeServerCreateDto.Name, Ip: nodeServerCreateDto.Ip}
	if err := service.CreateNodeServer(&nodeServer); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func SelectNodeServerPage(c *gin.Context) {
	var nodeServerPageDto dto.NodeServerPageDto
	_ = c.ShouldBindQuery(&nodeServerPageDto)
	if err := validate.Struct(&nodeServerPageDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	nodeServerPageVo, err := service.SelectNodeServerPage(nodeServerPageDto.Name, nodeServerPageDto.Ip, nodeServerPageDto.PageNum, nodeServerPageDto.PageSize, c)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nodeServerPageVo, c)
}

func DeleteNodeServerById(c *gin.Context) {
	var nodeServerRequireIdDto dto.RequiredIdDto
	_ = c.ShouldBindJSON(&nodeServerRequireIdDto)
	if err := validate.Struct(&nodeServerRequireIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	if err := service.DeleteNodeServerById(nodeServerRequireIdDto.Id); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func UpdateNodeServerById(c *gin.Context) {
	var nodeServerUpdateDto dto.NodeServerUpdateDto
	_ = c.ShouldBindJSON(&nodeServerUpdateDto)
	if err := validate.Struct(&nodeServerUpdateDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	if err := service.UpdateNodeServerById(&nodeServerUpdateDto); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func SelectNodeServerList(c *gin.Context) {
	var nodeServerDto dto.NodeServerDto
	_ = c.ShouldBindJSON(&nodeServerDto)
	if err := validate.Struct(&nodeServerDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	nodeServerListVos, err := service.SelectNodeServerList(&nodeServerDto)
	if err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	vo.Success(nodeServerListVos, c)
}
