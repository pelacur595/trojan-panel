package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"net/url"
	"strings"
	"sync"
	"trojan-panel/core"
	"trojan-panel/dao"
	"trojan-panel/module"
	"trojan-panel/module/constant"
	"trojan-panel/module/dto"
	"trojan-panel/module/vo"
	"trojan-panel/util"
)

func SelectNodeById(id *uint) (*vo.NodeOneVo, error) {
	node, err := dao.SelectNodeById(id)
	if err != nil {
		return nil, err
	}
	if node != nil {
		nodeOneVo := vo.NodeOneVo{
			Id:         *node.Id,
			NodeSubId:  *node.NodeSubId,
			NodeTypeId: *node.NodeTypeId,
			Name:       *node.Name,
			Ip:         *node.Ip,
			Port:       *node.Port,
			CreateTime: *node.CreateTime,
		}
		nodeTypeId := node.NodeTypeId
		switch *nodeTypeId {
		case 1:
			nodeXray, err := dao.SelectNodeXrayById(node.NodeSubId)
			if err != nil {
				return nil, err
			}
			nodeOneVo.XrayProtocol = nodeXray.Protocol
			xrayStreamSettingsEntity := vo.XrayStreamSettingsEntity{}
			if err = json.Unmarshal([]byte(*nodeXray.StreamSettings), &xrayStreamSettingsEntity); err != nil {
				logrus.Errorln(fmt.Sprintf("StreamSettings JSON反转失败 err: %v", err))
				return nil, errors.New(constant.SysError)
			}
			nodeOneVo.XrayStreamSettingsEntity = &xrayStreamSettingsEntity
		case 2:
			nodeTrojanGo, err := dao.SelectNodeTrojanGoById(node.NodeSubId)
			if err != nil {
				return nil, err
			}
			nodeOneVo.TrojanGoSni = nodeTrojanGo.Sni
			nodeOneVo.TrojanGoMuxEnable = nodeTrojanGo.MuxEnable
			nodeOneVo.TrojanGoWebsocketEnable = nodeTrojanGo.WebsocketEnable
			nodeOneVo.TrojanGoWebsocketPath = nodeTrojanGo.WebsocketPath
			nodeOneVo.TrojanGoWebsocketHost = nodeTrojanGo.WebsocketHost
			nodeOneVo.TrojanGoSsEnable = nodeTrojanGo.SsEnable
			nodeOneVo.TrojanGoSsMethod = nodeTrojanGo.SsMethod
			nodeOneVo.TrojanGoSsPassword = nodeTrojanGo.SsPassword
		case 3:
			nodeHysteria, err := dao.SelectNodeHysteriaById(node.NodeSubId)
			if err != nil {
				return nil, err
			}
			nodeOneVo.HysteriaProtocol = nodeHysteria.Protocol
			nodeOneVo.HysteriaUpMbps = nodeHysteria.UpMbps
			nodeOneVo.HysteriaDownMbps = nodeHysteria.DownMbps
		}
		return &nodeOneVo, nil
	}
	return nil, errors.New(constant.NodeNotExist)
}

func CreateNode(token string, nodeCreateDto dto.NodeCreateDto) error {
	// 校验端口
	var err error
	if nodeCreateDto.Port != nil && (*nodeCreateDto.Port <= 100 || *nodeCreateDto.Port >= 30000) {
		err = errors.New(constant.PortRangeError)
	}
	if *nodeCreateDto.NodeTypeId == 1 || *nodeCreateDto.NodeTypeId == 2 {
		if !util.IsPortAvailable(*nodeCreateDto.Port, "tcp") {
			err = errors.New(constant.PortIsOccupied)
		}
		if !util.IsPortAvailable(*nodeCreateDto.Port+10000, "tcp") {
			err = errors.New(constant.PortIsOccupied)
		}
	} else if *nodeCreateDto.NodeTypeId == 3 {
		if !util.IsPortAvailable(*nodeCreateDto.Port, "udp") {
			err = errors.New(constant.PortIsOccupied)
		}
	}
	if err != nil {
		return err
	}

	// 校验名称
	countName, err := dao.CountNodeByName(nodeCreateDto.Name)
	if err != nil {
		return err
	}
	if countName > 0 {
		return errors.New(constant.NodeNameExist)
	}

	var nodeId uint
	var mutex sync.Mutex
	defer mutex.Unlock()
	if mutex.TryLock() {
		// Grpc添加节点
		if err = GrpcAddNode(token, *nodeCreateDto.Ip, &core.NodeAddDto{
			NodeTypeId: uint64(*nodeCreateDto.NodeTypeId),
			//  Xray
			XrayPort:           uint64(*nodeCreateDto.Port),
			XrayProtocol:       *nodeCreateDto.XrayProtocol,
			XraySettings:       *nodeCreateDto.XraySettings,
			XrayStreamSettings: *nodeCreateDto.XrayStreamSettings,
			XrayTag:            *nodeCreateDto.XrayTag,
			XraySniffing:       *nodeCreateDto.XraySniffing,
			XrayAllocate:       *nodeCreateDto.XrayAllocate,
			// Trojan Go
			TrojanGoPort:            uint64(*nodeCreateDto.Port),
			TrojanGoIp:              *nodeCreateDto.Ip,
			TrojanGoSni:             *nodeCreateDto.TrojanGoSni,
			TrojanGoMuxEnable:       uint64(*nodeCreateDto.TrojanGoMuxEnable),
			TrojanGoWebsocketEnable: uint64(*nodeCreateDto.TrojanGoWebsocketEnable),
			TrojanGoWebsocketPath:   *nodeCreateDto.TrojanGoWebsocketPath,
			TrojanGoWebsocketHost:   *nodeCreateDto.TrojanGoWebsocketHost,
			TrojanGoSSEnable:        uint64(*nodeCreateDto.TrojanGoSsEnable),
			TrojanGoSSMethod:        *nodeCreateDto.TrojanGoSsMethod,
			TrojanGoSSPassword:      *nodeCreateDto.TrojanGoSsPassword,
			// Hysteria
			HysteriaPort:     uint64(*nodeCreateDto.Port),
			HysteriaProtocol: *nodeCreateDto.HysteriaProtocol,
			HysteriaIp:       *nodeCreateDto.Ip,
			HysteriaUpMbps:   int64(*nodeCreateDto.HysteriaUpMbps),
			HysteriaDownMbps: int64(*nodeCreateDto.HysteriaDownMbps),
		}); err != nil {
			return err
		}
		// 数据插入到数据库中
		if *nodeCreateDto.NodeTypeId == 1 {
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
		} else if *nodeCreateDto.NodeTypeId == 2 {
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
		} else if *nodeCreateDto.NodeTypeId == 3 {
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

		// 在主表中插入数据
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

func DeleteNodeById(token string, id *uint) error {
	var mutex sync.Mutex
	defer mutex.TryLock()
	if mutex.TryLock() {
		node, err := dao.SelectNodeById(id)
		if err != nil {
			return err
		}
		if err = GrpcRemoveNode(token, *node.Ip, *node.Port, *node.NodeTypeId); err != nil {
			return err
		}
		if err = dao.DeleteNodeById(id); err != nil {
			return err
		}
		if *node.NodeTypeId == 1 {
			if err := dao.DeleteNodeXrayById(node.NodeSubId); err != nil {
				return err
			}
		} else if *node.NodeTypeId == 2 {
			if err := dao.DeleteNodeTrojanGoById(node.NodeSubId); err != nil {
				return err
			}
		} else if *node.NodeTypeId == 3 {
			if err := dao.DeleteNodeHysteriaById(node.NodeSubId); err != nil {
				return err
			}
		}
	}
	return nil
}

func UpdateNodeById(token string, nodeUpdateDto *dto.NodeUpdateDto) error {
	count, err := dao.CountNodeByName(nodeUpdateDto.Name)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(constant.NodeNameExist)
	}
	var mutex sync.Mutex
	defer mutex.Unlock()
	if mutex.TryLock() {
		nodeEntity, err := dao.SelectNodeById(nodeUpdateDto.Id)
		if err != nil {
			return err
		}
		if nodeUpdateDto.NodeTypeId == nodeEntity.NodeTypeId {
			// Grpc的操作
			if err = GrpcRemoveNode(token, *nodeUpdateDto.Ip, *nodeEntity.Port, *nodeEntity.NodeTypeId); err != nil {
				return err
			}
			if err = GrpcAddNode(token, *nodeUpdateDto.Ip, &core.NodeAddDto{
				NodeTypeId: uint64(*nodeEntity.NodeTypeId),
				//  Xray
				XrayPort:           uint64(*nodeUpdateDto.Port),
				XrayProtocol:       *nodeUpdateDto.XrayProtocol,
				XraySettings:       *nodeUpdateDto.XraySettings,
				XrayStreamSettings: *nodeUpdateDto.XrayStreamSettings,
				XrayTag:            *nodeUpdateDto.XrayTag,
				XraySniffing:       *nodeUpdateDto.XraySniffing,
				XrayAllocate:       *nodeUpdateDto.XrayAllocate,
				// Trojan Go
				TrojanGoPort:            uint64(*nodeUpdateDto.Port),
				TrojanGoIp:              *nodeUpdateDto.Ip,
				TrojanGoSni:             *nodeUpdateDto.TrojanGoSni,
				TrojanGoMuxEnable:       uint64(*nodeUpdateDto.TrojanGoMuxEnable),
				TrojanGoWebsocketEnable: uint64(*nodeUpdateDto.TrojanGoWebsocketEnable),
				TrojanGoWebsocketPath:   *nodeUpdateDto.TrojanGoWebsocketPath,
				TrojanGoWebsocketHost:   *nodeUpdateDto.TrojanGoWebsocketHost,
				TrojanGoSSEnable:        uint64(*nodeUpdateDto.TrojanGoSsEnable),
				TrojanGoSSMethod:        *nodeUpdateDto.TrojanGoSsMethod,
				TrojanGoSSPassword:      *nodeUpdateDto.TrojanGoSsPassword,
				// Hysteria
				HysteriaPort:     uint64(*nodeUpdateDto.Port),
				HysteriaProtocol: *nodeUpdateDto.HysteriaProtocol,
				HysteriaIp:       *nodeUpdateDto.Ip,
				HysteriaUpMbps:   int64(*nodeUpdateDto.HysteriaUpMbps),
				HysteriaDownMbps: int64(*nodeUpdateDto.HysteriaDownMbps),
			}); err != nil {
				return err
			}
			// 修改了节点类型
			if *nodeEntity.NodeTypeId == 1 {
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
			} else if *nodeEntity.NodeTypeId == 2 {
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
			} else if *nodeEntity.NodeTypeId == 3 {
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
			node := module.Node{
				Id:   nodeUpdateDto.Id,
				Name: nodeUpdateDto.Name,
				Ip:   nodeUpdateDto.Ip,
				Port: nodeUpdateDto.Port,
			}
			if err = dao.UpdateNodeById(&node); err != nil {
				return err
			}
		} else {
			// 修改了节点类型需要删除分库的数据，然后重新再插入
			if err = DeleteNodeById(token, nodeUpdateDto.Id); err != nil {
				return err
			}
			if err = GrpcAddNode(token, *nodeUpdateDto.Ip, &core.NodeAddDto{
				NodeTypeId: uint64(*nodeUpdateDto.NodeTypeId),
				//  Xray
				XrayPort:           uint64(*nodeUpdateDto.Port),
				XrayProtocol:       *nodeUpdateDto.XrayProtocol,
				XraySettings:       *nodeUpdateDto.XraySettings,
				XrayStreamSettings: *nodeUpdateDto.XrayStreamSettings,
				XrayTag:            *nodeUpdateDto.XrayTag,
				XraySniffing:       *nodeUpdateDto.XraySniffing,
				XrayAllocate:       *nodeUpdateDto.XrayAllocate,
				// Trojan Go
				TrojanGoPort:            uint64(*nodeUpdateDto.Port),
				TrojanGoIp:              *nodeUpdateDto.Ip,
				TrojanGoSni:             *nodeUpdateDto.TrojanGoSni,
				TrojanGoMuxEnable:       uint64(*nodeUpdateDto.TrojanGoMuxEnable),
				TrojanGoWebsocketEnable: uint64(*nodeUpdateDto.TrojanGoWebsocketEnable),
				TrojanGoWebsocketPath:   *nodeUpdateDto.TrojanGoWebsocketPath,
				TrojanGoWebsocketHost:   *nodeUpdateDto.TrojanGoWebsocketHost,
				TrojanGoSSEnable:        uint64(*nodeUpdateDto.TrojanGoSsEnable),
				TrojanGoSSMethod:        *nodeUpdateDto.TrojanGoSsMethod,
				TrojanGoSSPassword:      *nodeUpdateDto.TrojanGoSsPassword,
				// Hysteria
				HysteriaPort:     uint64(*nodeUpdateDto.Port),
				HysteriaProtocol: *nodeUpdateDto.HysteriaProtocol,
				HysteriaIp:       *nodeUpdateDto.Ip,
				HysteriaUpMbps:   int64(*nodeUpdateDto.HysteriaUpMbps),
				HysteriaDownMbps: int64(*nodeUpdateDto.HysteriaDownMbps),
			}); err != nil {
				return err
			}
			var nodeId uint
			if *nodeUpdateDto.NodeTypeId == 1 {
				nodeXray := module.NodeXray{
					Protocol:       nodeUpdateDto.XrayProtocol,
					Settings:       nodeUpdateDto.XraySettings,
					StreamSettings: nodeUpdateDto.XrayStreamSettings,
					Tag:            nodeUpdateDto.XrayTag,
					Sniffing:       nodeUpdateDto.XraySniffing,
					Allocate:       nodeUpdateDto.XrayAllocate,
				}
				nodeId, err = dao.CreateNodeXray(&nodeXray)
				if err != nil {
					return nil
				}
			} else if *nodeUpdateDto.NodeTypeId == 2 {
				trojanGo := module.NodeTrojanGo{
					Sni:             nodeUpdateDto.TrojanGoSni,
					MuxEnable:       nodeUpdateDto.TrojanGoMuxEnable,
					WebsocketEnable: nodeUpdateDto.TrojanGoWebsocketEnable,
					WebsocketPath:   nodeUpdateDto.TrojanGoWebsocketPath,
					WebsocketHost:   nodeUpdateDto.TrojanGoWebsocketHost,
					SsEnable:        nodeUpdateDto.TrojanGoSsEnable,
					SsMethod:        nodeUpdateDto.TrojanGoSsMethod,
					SsPassword:      nodeUpdateDto.TrojanGoSsPassword,
				}
				nodeId, err = dao.CreateNodeTrojanGo(&trojanGo)
				if err != nil {
					return nil
				}
			} else if *nodeUpdateDto.NodeTypeId == 3 {
				hysteria := module.NodeHysteria{
					Protocol: nodeUpdateDto.HysteriaProtocol,
					UpMbps:   nodeUpdateDto.HysteriaUpMbps,
					DownMbps: nodeUpdateDto.HysteriaDownMbps,
				}
				nodeId, err = dao.CreateNodeHysteria(&hysteria)
				if err != nil {
					return nil
				}
			}

			node := module.Node{
				Id:         nodeUpdateDto.Id,
				NodeSubId:  &nodeId,
				NodeTypeId: nodeUpdateDto.NodeTypeId,
				Name:       nodeUpdateDto.Name,
				Ip:         nodeUpdateDto.Ip,
				Port:       nodeUpdateDto.Port,
			}
			if err = dao.UpdateNodeById(&node); err != nil {
				return err
			}
		}
	}
	return nil
}

func NodeQRCode(accountId *uint, id *uint) ([]byte, error) {
	nodeUrl, err := NodeURL(accountId, id)
	if err != nil {
		return nil, err
	}
	// 生成二维码
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

	password, err := dao.SelectConnectPassword(accountId, nil)
	if err != nil {
		return "", errors.New(constant.NodeURLError)
	}

	// 构建URL
	var headBuilder strings.Builder

	if *nodeType.Name == constant.XrayName {
		// todo
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
		nodeHysteria, err := dao.SelectNodeHysteriaById(node.NodeSubId)
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

func GrpcAddNode(token string, ip string, nodeAddDto *core.NodeAddDto) error {
	if err := core.AddNode(token, ip, nodeAddDto); err != nil {
		logrus.Errorf("gRPC添加节点异常 ip: %s err: %v", ip, err)
		return errors.New(constant.GrpcAddNodeError)
	}
	return nil
}

func GrpcRemoveNode(token string, ip string, port uint, nodeType uint) error {
	if err := core.RemoveNode(token, ip, &core.NodeRemoveDto{
		NodeType: uint64(nodeType),
		Port:     uint64(port),
	}); err != nil {
		logrus.Errorf("gRPC移除节点异常 ip: %s err: %v", ip, err)
	}
	return nil
}
