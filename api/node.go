package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"trojan/module"
	"trojan/module/constant"
	"trojan/module/dto"
	"trojan/module/vo"
	"trojan/service"
)

func SelectNodeById(c *gin.Context) {
	var nodeRequireIdDto dto.RequiredIdDto
	_ = c.ShouldBindQuery(&nodeRequireIdDto)
	if err := validate.Struct(&nodeRequireIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	nodeVo, err := service.SelectNodeById(nodeRequireIdDto.Id)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nodeVo, c)
}

func CreateNode(c *gin.Context) {
	var nodeCreateDto dto.NodeCreateDto
	_ = c.ShouldBindJSON(&nodeCreateDto)
	if err := validate.Struct(&nodeCreateDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	node := module.Node{
		Name:             nodeCreateDto.Name,
		Ip:               nodeCreateDto.Ip,
		Port:             nodeCreateDto.Port,
		Sni:              nodeCreateDto.Sni,
		Type:             nodeCreateDto.Type,
		WebsocketEnable:  nodeCreateDto.WebsocketEnable,
		WebsocketPath:    nodeCreateDto.WebsocketPath,
		SsEnable:         nodeCreateDto.SsEnable,
		SsMethod:         nodeCreateDto.SsMethod,
		SsPassword:       nodeCreateDto.SsPassword,
		HysteriaProtocol: nodeCreateDto.HysteriaProtocol,
		HysteriaUpMbps:   nodeCreateDto.HysteriaUpMbps,
		HysteriaDownMbps: nodeCreateDto.HysteriaDownMbps,
	}
	if err := service.CreateNode(&node); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func SelectNodePage(c *gin.Context) {
	var nodePageDto dto.NodePageDto
	_ = c.ShouldBindQuery(&nodePageDto)
	if err := validate.Struct(&nodePageDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	nodePageVo, err := service.SelectNodePage(nodePageDto.Name, nodePageDto.PageNum, nodePageDto.PageSize)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nodePageVo, c)
}

func DeleteNodeById(c *gin.Context) {
	var nodeRequireIdDto dto.RequiredIdDto
	_ = c.ShouldBindJSON(&nodeRequireIdDto)
	if err := validate.Struct(&nodeRequireIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	if err := service.DeleteNodeById(nodeRequireIdDto.Id); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func UpdateNodeById(c *gin.Context) {
	var nodeUpdateDto dto.NodeUpdateDto
	_ = c.ShouldBindJSON(&nodeUpdateDto)
	if err := validate.Struct(&nodeUpdateDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	node := module.Node{
		Id:               nodeUpdateDto.Id,
		Name:             nodeUpdateDto.Name,
		Ip:               nodeUpdateDto.Ip,
		Sni:              nodeUpdateDto.Sni,
		Port:             nodeUpdateDto.Port,
		Type:             nodeUpdateDto.Type,
		WebsocketEnable:  nodeUpdateDto.WebsocketEnable,
		WebsocketPath:    nodeUpdateDto.WebsocketPath,
		SsEnable:         nodeUpdateDto.SsEnable,
		SsMethod:         nodeUpdateDto.SsMethod,
		SsPassword:       nodeUpdateDto.SsPassword,
		HysteriaProtocol: nodeUpdateDto.HysteriaProtocol,
		HysteriaUpMbps:   nodeUpdateDto.HysteriaUpMbps,
		HysteriaDownMbps: nodeUpdateDto.HysteriaDownMbps,
	}
	if err := service.UpdateNodeById(&node); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func NodeQRCode(c *gin.Context) {
	var nodeQRCodeDto dto.NodeQRCodeDto
	_ = c.ShouldBindJSON(&nodeQRCodeDto)
	if err := validate.Struct(&nodeQRCodeDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	userInfo, err := service.GetUserInfo(c)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	qrCode, err := service.NodeQRCode(&userInfo.Id, nodeQRCodeDto)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(qrCode, c)
}

func NodeURL(c *gin.Context) {
	var nodeQRCodeDto dto.NodeQRCodeDto
	_ = c.ShouldBindJSON(&nodeQRCodeDto)
	if err := validate.Struct(&nodeQRCodeDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		fmt.Println(err.Error())
		return
	}
	userInfo, err := service.GetUserInfo(c)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	url, err := service.NodeURL(&userInfo.Id, nodeQRCodeDto)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(url, c)
}
