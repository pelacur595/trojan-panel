package service

import (
	"errors"
	"fmt"
	"github.com/skip2/go-qrcode"
	"net/url"
	"strings"
	"trojan/dao"
	"trojan/module"
	"trojan/module/constant"
	"trojan/module/dto"
	"trojan/module/vo"
)

func SelectNodeById(id *uint) (*vo.NodeVo, error) {
	nodeVo, err := dao.SelectNodeById(id)
	if err != nil {
		return nil, err
	}
	return nodeVo, nil
}

func CreateNode(node *module.Node) error {
	count, err := dao.CountNodeByName(node.Name)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(constant.NodeNameExist)
	}
	if err := dao.CreateNode(node); err != nil {
		return err
	}
	return nil
}

func SelectNodePage(queryName *string, pageNum *uint, pageSize *uint) (*vo.NodePageVo, error) {
	nodePagVo, err := dao.SelectNodePage(queryName, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	return nodePagVo, nil
}

func DeleteNodeById(id *uint) error {
	if err := dao.DeleteNodeById(id); err != nil {
		return err
	}
	return nil
}

func UpdateNodeById(node *module.Node) error {
	if err := dao.UpdateNodeById(node); err != nil {
		return err
	}
	return nil
}

func NodeQRCode(userId *uint, NodeQRCodeDto dto.NodeQRCodeDto) ([]byte, error) {
	nodeUrl, err := NodeURL(userId, NodeQRCodeDto)
	if err != nil {
		return nil, err
	}
	// Shadowrocket暂不支持trojan-go
	nodeUrl = strings.Replace(nodeUrl, "trojan-go", "trojan", 1)
	qrCode, err := qrcode.Encode(nodeUrl, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	return qrCode, nil
}

func NodeURL(userId *uint, nodeQRCodeDto dto.NodeQRCodeDto) (string, error) {
	password, err := dao.UserQRCode(userId)
	if err != nil {
		return "", err
	}
	nodeTypeVo, err := dao.SelectNodeTypeById(nodeQRCodeDto.Type)
	if err != nil {
		return "", err
	}
	if nodeTypeVo.Prefix == "" ||
		nodeQRCodeDto.Ip == nil || *nodeQRCodeDto.Ip == "" ||
		nodeQRCodeDto.Port == nil || *nodeQRCodeDto.Port == 0 {
		return "", errors.New(constant.NodeURLError)
	}
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s://%s@%s:%d", nodeTypeVo.Prefix, password, *nodeQRCodeDto.Ip, *nodeQRCodeDto.Port))
	if nodeTypeVo.Prefix == constant.TrojanGoPrefix {
		if nodeQRCodeDto.WebsocketEnable != nil && *nodeQRCodeDto.WebsocketEnable != 0 &&
			nodeQRCodeDto.WebsocketPath != nil && *nodeQRCodeDto.WebsocketPath != "" {
			builder.WriteString("&type=ws")
			builder.WriteString(fmt.Sprintf("&path=%s", *nodeQRCodeDto.WebsocketPath))
		}
		if nodeQRCodeDto.SsEnable != nil && *nodeQRCodeDto.SsEnable != 0 ||
			nodeQRCodeDto.SsMethod != nil && *nodeQRCodeDto.SsMethod != "" ||
			nodeQRCodeDto.SsPassword != nil && *nodeQRCodeDto.SsPassword != "" {
			builder.WriteString(fmt.Sprintf("encryption=ss;%s:%s", *nodeQRCodeDto.SsMethod, *nodeQRCodeDto.SsPassword))
		}
	}
	if nodeQRCodeDto.Name != nil && *nodeQRCodeDto.Name != "" {
		builder.WriteString(fmt.Sprintf("#%s", *nodeQRCodeDto.Name))
	}
	return url.PathEscape(builder.String()), nil
}

func CountNode() (int, error) {
	nodeCount, err := dao.CountNode()
	if err != nil {
		return 0, err
	}
	return nodeCount, nil
}
