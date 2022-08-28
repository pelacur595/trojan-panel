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
	"trojan/util"
)

func SelectNodeById(id *uint) (*module.Node, error) {
	return dao.SelectNodeById(id)
}

func CreateNode(nodeCreateDto dto.NodeCreateDto) error {
	count, err := dao.CountNodeByName(nodeCreateDto.Name)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(constant.NodeNameExist)
	}
	nodeType, err := dao.SelectNodeTypeById(nodeCreateDto.NodeTypeId)
	if err != nil {
		return err
	}

	var nodeId uint
	if nodeType.Name == constant.TrojanGoName {
		trojanGo := module.NodeTrojanGo{
			Sni:             nodeCreateDto.TrojanGoSni,
			MuxEnable:       nodeCreateDto.TrojanGoMuxEnable,
			WebsocketEnable: nodeCreateDto.TrojanGoWebsocketEnable,
			WebsocketPath:   nodeCreateDto.TrojanGoWebsocketPath,
			SsEnable:        nodeCreateDto.TrojanGoSsEnable,
			SsMethod:        nodeCreateDto.TrojanGoSsMethod,
			SsPassword:      nodeCreateDto.TrojanGoSsPassword,
		}
		nodeId, err = dao.CreateNodeTrojanGo(&trojanGo)
		if err != nil {
			return nil
		}
	}

	if nodeType.Name == constant.HysteriaName {
		hysteria := module.NodeHysteria{
			Protocol: nodeCreateDto.HysteriaProtocol,
			UpMbps:   nodeCreateDto.HysteriaUpMbps,
			DownMbps: nodeCreateDto.HysteriaDownMbps,
		}
		nodeId, err = dao.CreateNodeHysteria(&hysteria)
		if err != nil {
			return nil
		}
	}

	if nodeType.Name == constant.XrayName {
		nodeXray := module.NodeXray{
			Protocol:       nodeCreateDto.XrayProtocol,
			Settings:       nodeCreateDto.XraySettings,
			StreamSettings: nodeCreateDto.XrayStreamSettings,
			Tag:            nodeCreateDto.XrayTag,
			Sniffing:       nodeCreateDto.XraySniffing,
			Allocate:       nodeCreateDto.XrayAllocate,
		}
		nodeId, err = dao.CreateNodeXray(&nodeXray)
		if err != nil {
			return nil
		}
	}
	node := module.Node{
		NodeSubId:  &nodeId,
		NodeTypeId: nodeCreateDto.NodeTypeId,
		Name:       nodeCreateDto.Name,
		Ip:         nodeCreateDto.Ip,
		Port:       nodeCreateDto.Port,
	}
	if err = dao.CreateNode(&node); err != nil {
		return err
	}
	return nil
}

func SelectNodePage(queryName *string, pageNum *uint, pageSize *uint) (*vo.NodePageVo, error) {
	nodePagVo, err := dao.SelectNodePage(queryName, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	for index, node := range nodePagVo.Nodes {
		ttl, err := util.Ping(node.Ip)
		if err != nil {
			nodePagVo.Nodes[index].Ping = -1
		}
		nodePagVo.Nodes[index].Ping = ttl
	}
	return nodePagVo, nil
}

func DeleteNodeById(id *uint) error {
	return dao.DeleteNodeById(id)
}

func UpdateNodeById(nodeUpdateDto *dto.NodeUpdateDto) error {

	node := module.Node{
		Id:   nodeUpdateDto.Id,
		Name: nodeUpdateDto.Name,
		Ip:   nodeUpdateDto.Ip,
		Port: nodeUpdateDto.Port,
	}
	if err := dao.UpdateNodeById(&node); err != nil {
		return err
	}

	nodeType, err := dao.SelectNodeTypeById(nodeUpdateDto.NodeTypeId)
	if err != nil {
		return nil
	}
	if nodeType.Name == constant.TrojanGoName {
		nodeTrojanGo := module.NodeTrojanGo{
			Id:              nodeUpdateDto.NodeSubId,
			Sni:             nodeUpdateDto.TrojanGoSni,
			MuxEnable:       nodeUpdateDto.TrojanGoMuxEnable,
			WebsocketEnable: nodeUpdateDto.TrojanGoWebsocketEnable,
			WebsocketPath:   nodeUpdateDto.TrojanGoWebsocketPath,
			SsEnable:        nodeUpdateDto.TrojanGoSsEnable,
			SsMethod:        nodeUpdateDto.TrojanGoSsMethod,
			SsPassword:      nodeUpdateDto.TrojanGoSsPassword,
		}
		if err := dao.UpdateNodeTrojanGoById(&nodeTrojanGo); err != nil {
			return err
		}
	}
	if nodeType.Name == constant.HysteriaName {
		nodeHysteria := module.NodeHysteria{
			Id:       nodeUpdateDto.NodeSubId,
			Protocol: nodeUpdateDto.HysteriaProtocol,
			UpMbps:   nodeUpdateDto.HysteriaUpMbps,
			DownMbps: nodeUpdateDto.HysteriaDownMbps,
		}
		if err := dao.UpdateNodeHysteriaById(&nodeHysteria); err != nil {
			return err
		}
	}
	if nodeType.Name == constant.XrayName {
		nodeXray := module.NodeXray{
			Id:             nodeUpdateDto.NodeSubId,
			Protocol:       nodeUpdateDto.XrayProtocol,
			Settings:       nodeUpdateDto.XraySettings,
			StreamSettings: nodeUpdateDto.XrayStreamSettings,
			Tag:            nodeUpdateDto.XrayTag,
			Sniffing:       nodeUpdateDto.XraySniffing,
			Allocate:       nodeUpdateDto.XrayAllocate,
		}
		if err := dao.UpdateNodeXrayById(&nodeXray); err != nil {
			return err
		}
	}
	return nil
}

func NodeQRCode(userId *uint, id *uint) ([]byte, error) {
	nodeUrl, err := NodeURL(userId, id)
	if err != nil {
		return nil, err
	}
	qrCode, err := qrcode.Encode(nodeUrl, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	return qrCode, nil
}

func NodeURL(userId *uint, id *uint) (string, error) {

	node, err := dao.SelectNodeById(id)
	if err != nil {
		return "", nil
	}

	nodeType, err := dao.SelectNodeTypeById(node.NodeTypeId)
	if err != nil {
		return "", err
	}

	password, err := dao.AccountQRCode(userId)
	if err != nil {
		return "", err
	}

	// 构建URL
	var headBuilder strings.Builder
	if nodeType.Name == constant.TrojanGoName {
		nodeTrojanGo, err := dao.SelectNodeTrojanGoById(node.NodeSubId)
		if err != nil {
			return "", nil
		}
		headBuilder.WriteString(fmt.Sprintf("trojan-go://%s@%s:%d?", url.PathEscape(password),
			*node.Ip, *node.Port))
		var sni string
		if nodeTrojanGo.Sni != nil && *nodeTrojanGo.Sni != "" {
			sni = *nodeTrojanGo.Sni
		} else {
			sni = *node.Ip
		}
		headBuilder.WriteString(fmt.Sprintf("sni=%s", url.PathEscape(sni)))
		if nodeTrojanGo.WebsocketEnable != nil && *nodeTrojanGo.WebsocketEnable != 0 &&
			nodeTrojanGo.WebsocketPath != nil && *nodeTrojanGo.WebsocketPath != "" {
			headBuilder.WriteString(fmt.Sprintf("&type=%s", url.PathEscape("ws")))
			headBuilder.WriteString(fmt.Sprintf("&path=%s",
				url.PathEscape(fmt.Sprintf("/%s", *nodeTrojanGo.WebsocketPath))))
			if nodeTrojanGo.SsEnable != nil && *nodeTrojanGo.SsEnable != 0 ||
				nodeTrojanGo.SsMethod != nil && *nodeTrojanGo.SsMethod != "" ||
				nodeTrojanGo.SsPassword != nil && *nodeTrojanGo.SsPassword != "" {
				headBuilder.WriteString(fmt.Sprintf("&encryption=%s", url.PathEscape(
					fmt.Sprintf("ss;%s:%s", *nodeTrojanGo.SsMethod, *nodeTrojanGo.SsPassword))))
			}
		}
	}

	if nodeType.Name == constant.HysteriaName {
		nodeHysteria, err := dao.SelectHysteriaById(id)
		if err != nil {
			return "", err
		}
		headBuilder.WriteString(fmt.Sprintf("hysteria://%s:%d?protocol=%s&auth=%s&upmbps=%d&downmbps=%d",
			node.Ip,
			node.Port,
			*nodeHysteria.Protocol,
			password,
			*nodeHysteria.UpMbps,
			*nodeHysteria.DownMbps))
	}

	if nodeType.Name == constant.XrayName {

	}

	if node.Name != nil && *node.Name != "" {
		headBuilder.WriteString(fmt.Sprintf("#%s", url.PathEscape(*node.Name)))
	}
	return headBuilder.String(), nil
}

func CountNode() (int, error) {
	return dao.CountNode()
}
