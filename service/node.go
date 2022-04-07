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

func NodeQRCode(userId *uint, name *string, ip *string, port *uint, nodeType *uint) ([]byte, error) {
	nodeUrl, err := NodeURL(userId, name, ip, port, nodeType)
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

func NodeURL(userId *uint, name *string, ip *string, port *uint, nodeType *uint) (string, error) {
	password, err := dao.UserQRCode(userId)
	if err != nil {
		return "", err
	}
	nodeTypeVo, err := dao.SelectNodeTypeById(nodeType)
	if err != nil {
		return "", err
	}
	if nodeTypeVo.Prefix == "" || ip == nil || *ip == "" || port == nil || *port == 0 {
		return "", errors.New(constant.NodeURLError)
	}
	return fmt.Sprintf("%s://%s@%s:%d#%s", nodeTypeVo.Prefix, password, *ip, *port, url.PathEscape(*name)), nil
}

func CountNode() (int, error) {
	nodeCount, err := dao.CountNode()
	if err != nil {
		return 0, err
	}
	return nodeCount, nil
}

func CountNodeByName(queryName *string) (int, error) {
	nodeCount, err := dao.CountNodeByName(queryName)
	if err != nil {
		return 0, err
	}
	return nodeCount, nil
}
