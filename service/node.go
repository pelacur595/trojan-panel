package service

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"net/url"
	"strings"
	"sync"
	"trojan/core"
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
	var mutex sync.Mutex
	defer mutex.Unlock()
	if mutex.TryLock() {
		//if err = GrpcAddNode(&core.NodeAddDto{
		//	NodeType:                uint64(*nodeType.Id),
		//	TrojanGoPort:            uint64(*nodeCreateDto.Port),
		//	TrojanGoIp:              *nodeCreateDto.Ip,
		//	TrojanGoSni:             *nodeCreateDto.TrojanGoSni,
		//	TrojanGoMuxEnable:       uint64(*nodeCreateDto.TrojanGoMuxEnable),
		//	TrojanGoWebsocketEnable: uint64(*nodeCreateDto.TrojanGoWebsocketEnable),
		//	TrojanGoWebsocketPath:   *nodeCreateDto.TrojanGoWebsocketPath,
		//	TrojanGoWebsocketHost:   *nodeCreateDto.TrojanGoWebsocketHost,
		//	TrojanGoSSEnable:        uint64(*nodeCreateDto.TrojanGoSsEnable),
		//	TrojanGoSSMethod:        *nodeCreateDto.TrojanGoSsMethod,
		//	TrojanGoSSPassword:      *nodeCreateDto.TrojanGoSsPassword,
		//	HysteriaPort:            uint64(*nodeCreateDto.Port),
		//	HysteriaProtocol:        *nodeCreateDto.HysteriaProtocol,
		//	HysteriaIp:              *nodeCreateDto.Ip,
		//	HysteriaUpMbps:          int64(*nodeCreateDto.HysteriaUpMbps),
		//	HysteriaDownMbps:        int64(*nodeCreateDto.HysteriaDownMbps),
		//	XrayPort:                uint64(*nodeCreateDto.Port),
		//	XrayProtocol:            *nodeCreateDto.XrayProtocol,
		//	XraySettings:            *nodeCreateDto.XraySettings,
		//	XrayStreamSettings:      *nodeCreateDto.XrayStreamSettings,
		//	XrayTag:                 *nodeCreateDto.XrayTag,
		//	XraySniffing:            *nodeCreateDto.XraySniffing,
		//	XrayAllocate:            *nodeCreateDto.XrayAllocate,
		//}); err != nil {
		//	return err
		//}
		if *nodeType.Name == constant.TrojanGoName {
			trojanGo := module.NodeTrojanGo{
				Sni:             nodeCreateDto.TrojanGoSni,
				MuxEnable:       nodeCreateDto.TrojanGoMuxEnable,
				WebsocketEnable: nodeCreateDto.TrojanGoWebsocketEnable,
				WebsocketPath:   nodeCreateDto.TrojanGoWebsocketPath,
				WebsocketHost:   nodeCreateDto.TrojanGoWebsocketHost,
				SsEnable:        nodeCreateDto.TrojanGoSsEnable,
				SsMethod:        nodeCreateDto.TrojanGoSsMethod,
				SsPassword:      nodeCreateDto.TrojanGoSsPassword,
			}
			nodeId, err = dao.CreateNodeTrojanGo(&trojanGo)
			if err != nil {
				return nil
			}
		} else if *nodeType.Name == constant.HysteriaName {
			hysteria := module.NodeHysteria{
				Protocol: nodeCreateDto.HysteriaProtocol,
				UpMbps:   nodeCreateDto.HysteriaUpMbps,
				DownMbps: nodeCreateDto.HysteriaDownMbps,
			}
			nodeId, err = dao.CreateNodeHysteria(&hysteria)
			if err != nil {
				return nil
			}
		} else if *nodeType.Name == constant.XrayName {
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
	nodePage, total, err := dao.SelectNodePage(queryName, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	nodeVos := make([]vo.NodeVo, 0)
	for _, item := range *nodePage {
		nodeVo := vo.NodeVo{
			Id:         *item.Id,
			NodeSubId:  *item.NodeSubId,
			NodeTypeId: *item.NodeTypeId,
			Name:       *item.Name,
			Ip:         *item.Ip,
			Port:       *item.Port,
			CreateTime: *item.CreateTime,
		}

		ttl, err := util.Ping(*item.Ip)
		if err != nil {
			ttl = -1
		}
		nodeVo.Ping = ttl

		//nodeType, err := dao.SelectNodeTypeById(item.NodeTypeId)
		//if err != nil {
		//	continue
		//}
		//if *nodeType.Name == constant.TrojanGoName {
		//	nodeTroajanGo, err := dao.SelectNodeTrojanGoById(item.NodeSubId)
		//	if err != nil {
		//		continue
		//	}
		//	nodeVo.TrojanGoSni = nodeTroajanGo.Sni
		//	nodeVo.TrojanGoMuxEnable = nodeTroajanGo.MuxEnable
		//	nodeVo.TrojanGoWebsocketEnable = nodeTroajanGo.WebsocketEnable
		//	nodeVo.TrojanGoWebsocketPath = nodeTroajanGo.WebsocketPath
		//	nodeVo.TrojanGoWebsocketHost = nodeTroajanGo.WebsocketHost
		//	nodeVo.TrojanGoSsEnable = nodeTroajanGo.SsEnable
		//	nodeVo.TrojanGoSsMethod = nodeTroajanGo.SsMethod
		//	nodeVo.TrojanGoSsPassword = nodeTroajanGo.SsPassword
		//}
		//if *nodeType.Name == constant.HysteriaName {
		//	nodeHysteria, err := dao.SelectNodeHysteriaById(item.NodeSubId)
		//	if err != nil {
		//		continue
		//	}
		//	nodeVo.HysteriaProtocol = nodeHysteria.Protocol
		//	nodeVo.HysteriaUpMbps = nodeHysteria.UpMbps
		//	nodeVo.HysteriaDownMbps = nodeHysteria.DownMbps
		//}
		//if *nodeType.Name == constant.XrayName {
		//	nodeXray, err := dao.SelectNodeXrayById(item.NodeSubId)
		//	if err != nil {
		//		continue
		//	}
		//	nodeVo.XrayProtocol = nodeXray.Protocol
		//	nodeVo.XraySettings = nodeXray.Settings
		//	nodeVo.XrayStreamSettings = nodeXray.StreamSettings
		//	nodeVo.XrayTag = nodeXray.Tag
		//	nodeVo.XraySniffing = nodeXray.Sniffing
		//	nodeVo.XrayAllocate = nodeXray.Allocate
		//}

		nodeVos = append(nodeVos, nodeVo)
	}
	nodePageVo := vo.NodePageVo{
		BaseVoPage: vo.BaseVoPage{
			PageNum:  *pageNum,
			PageSize: *pageSize,
			Total:    total,
		},
		Nodes: nodeVos,
	}
	return &nodePageVo, nil
}

func DeleteNodeById(id *uint) error {
	var mutex sync.Mutex
	defer mutex.TryLock()
	if mutex.TryLock() {
		node, err := dao.SelectNodeById(id)
		if err != nil {
			return err
		}
		if err = GrpcRemoveNode(*node.NodeTypeId, *node.Port); err != nil {
			return err
		}
		if err = dao.DeleteNodeById(id); err != nil {
			return err
		}
	}
	return nil
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
	if *nodeType.Name == constant.XrayName {
		nodeXray := module.NodeXray{
			Id:             nodeUpdateDto.NodeSubId,
			Protocol:       nodeUpdateDto.XrayProtocol,
			Settings:       nodeUpdateDto.XraySettings,
			StreamSettings: nodeUpdateDto.XrayStreamSettings,
			Tag:            nodeUpdateDto.XrayTag,
			Sniffing:       nodeUpdateDto.XraySniffing,
			Allocate:       nodeUpdateDto.XrayAllocate,
		}
		if err = dao.UpdateNodeXrayById(&nodeXray); err != nil {
			return err
		}
	} else if *nodeType.Name == constant.TrojanGoName {
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
		if err = dao.UpdateNodeTrojanGoById(&nodeTrojanGo); err != nil {
			return err
		}
	} else if *nodeType.Name == constant.HysteriaName {
		nodeHysteria := module.NodeHysteria{
			Id:       nodeUpdateDto.NodeSubId,
			Protocol: nodeUpdateDto.HysteriaProtocol,
			UpMbps:   nodeUpdateDto.HysteriaUpMbps,
			DownMbps: nodeUpdateDto.HysteriaDownMbps,
		}
		if err = dao.UpdateNodeHysteriaById(&nodeHysteria); err != nil {
			return err
		}
	}
	return nil
}

func NodeQRCode(accountId *uint, id *uint) ([]byte, error) {
	nodeUrl, err := NodeURL(accountId, id)
	if err != nil {
		return nil, err
	}
	qrCode, err := qrcode.Encode(nodeUrl, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	return qrCode, nil
}

func NodeURL(accountId *uint, id *uint) (string, error) {

	node, err := dao.SelectNodeById(id)
	if err != nil {
		return "", errors.New(constant.NodeURLError)
	}

	nodeType, err := dao.SelectNodeTypeById(node.NodeTypeId)
	if err != nil {
		return "", errors.New(constant.NodeURLError)
	}

	password, err := dao.AccountQRCode(accountId)
	if err != nil {
		return "", errors.New(constant.NodeURLError)
	}

	// 构建URL
	var headBuilder strings.Builder

	if *nodeType.Name == constant.XrayName {

	} else if *nodeType.Name == constant.TrojanGoName {
		nodeTrojanGo, err := dao.SelectNodeTrojanGoById(node.NodeSubId)
		if err != nil {
			return "", errors.New(constant.NodeURLError)
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
	} else if *nodeType.Name == constant.HysteriaName {
		nodeHysteria, err := dao.SelectNodeHysteriaById(id)
		if err != nil {
			return "", errors.New(constant.NodeURLError)
		}
		headBuilder.WriteString(fmt.Sprintf("hysteria://%s:%d?protocol=%s&auth=%s&upmbps=%d&downmbps=%d",
			*node.Ip,
			*node.Port,
			*nodeHysteria.Protocol,
			password,
			*nodeHysteria.UpMbps,
			*nodeHysteria.DownMbps))
	}

	if node.Name != nil && *node.Name != "" {
		headBuilder.WriteString(fmt.Sprintf("#%s", url.PathEscape(*node.Name)))
	}
	return headBuilder.String(), nil
}

func CountNode() (int, error) {
	return dao.CountNode()
}

func GrpcAddNode(nodeAddDto *core.NodeAddDto) error {
	nodes, err := dao.SelectNodesIpAndPort()
	if err != nil {
		return err
	}
	for _, node := range nodes {
		if err = core.AddNode(*node.Ip, nodeAddDto); err != nil {
			logrus.Errorf("gRPC添加节点异常 ip: %s err: %v", *node.Ip, err)
		}
		continue
	}
	return nil
}

func GrpcRemoveNode(nodeType uint, port uint) error {
	nodes, err := dao.SelectNodesIpAndPort()
	if err != nil {
		return err
	}
	for _, node := range nodes {
		if err = core.RemoveNode(*node.Ip, &core.NodeRemoveDto{
			NodeType: uint64(nodeType),
			Port:     uint64(port),
		}); err != nil {
			logrus.Errorf("gRPC添加节点异常 ip: %s err: %v", *node.Ip, err)
		}
		continue
	}
	return nil
}
