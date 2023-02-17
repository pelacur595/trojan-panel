package service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"net/url"
	"strings"
	"sync"
	"time"
	"trojan-panel/core"
	"trojan-panel/dao"
	"trojan-panel/module"
	"trojan-panel/module/bo"
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
			Id:           *node.Id,
			NodeServerId: *node.NodeServerId,
			NodeSubId:    *node.NodeSubId,
			NodeTypeId:   *node.NodeTypeId,
			Name:         *node.Name,
			Domain:       *node.Domain,
			Port:         *node.Port,
			CreateTime:   *node.CreateTime,
		}
		nodeTypeId := node.NodeTypeId
		switch *nodeTypeId {
		case constant.Xray:
			nodeXray, err := dao.SelectNodeXrayById(node.NodeSubId)
			if err != nil {
				return nil, err
			}
			nodeOneVo.XrayProtocol = *nodeXray.Protocol
			nodeOneVo.XrayFlow = *nodeXray.XrayFlow
			nodeOneVo.XraySSMethod = *nodeXray.XraySSMethod
			xrayStreamSettingsEntity := vo.XrayStreamSettingsEntity{}
			if nodeXray.StreamSettings != nil && *nodeXray.StreamSettings != "" {
				if err = json.Unmarshal([]byte(*nodeXray.StreamSettings), &xrayStreamSettingsEntity); err != nil {
					logrus.Errorln(fmt.Sprintf("StreamSettings JSON反转失败 err: %v", err))
					return nil, errors.New(constant.SysError)
				}
			}
			nodeOneVo.XrayStreamSettingsEntity = xrayStreamSettingsEntity
			nodeOneVo.XrayTag = *nodeXray.Tag
		case constant.TrojanGo:
			nodeTrojanGo, err := dao.SelectNodeTrojanGoById(node.NodeSubId)
			if err != nil {
				return nil, err
			}
			nodeOneVo.TrojanGoSni = *nodeTrojanGo.Sni
			nodeOneVo.TrojanGoMuxEnable = *nodeTrojanGo.MuxEnable
			nodeOneVo.TrojanGoWebsocketEnable = *nodeTrojanGo.WebsocketEnable
			nodeOneVo.TrojanGoWebsocketPath = *nodeTrojanGo.WebsocketPath
			nodeOneVo.TrojanGoWebsocketHost = *nodeTrojanGo.WebsocketHost
			nodeOneVo.TrojanGoSsEnable = *nodeTrojanGo.SsEnable
			nodeOneVo.TrojanGoSsMethod = *nodeTrojanGo.SsMethod
			nodeOneVo.TrojanGoSsPassword = *nodeTrojanGo.SsPassword
		case constant.Hysteria:
			nodeHysteria, err := dao.SelectNodeHysteriaById(node.NodeSubId)
			if err != nil {
				return nil, err
			}
			nodeOneVo.HysteriaProtocol = *nodeHysteria.Protocol
			nodeOneVo.HysteriaUpMbps = *nodeHysteria.UpMbps
			nodeOneVo.HysteriaDownMbps = *nodeHysteria.DownMbps
		}
		return &nodeOneVo, nil
	}
	return nil, errors.New(constant.NodeNotExist)
}
func SelectNodeInfo(id *uint, c *gin.Context) (*vo.NodeOneVo, error) {

	nodeOneVo, err := SelectNodeById(id)
	if err != nil {
		return nil, err
	}
	accountInfo := util.GetCurrentAccount(c)
	account, err := dao.SelectAccountByUsername(&accountInfo.Username)
	if err != nil {
		return nil, err
	}
	nodeOneVo.Password = *account.Pass
	if nodeOneVo.NodeTypeId == 1 && (nodeOneVo.XrayProtocol == "vless" || nodeOneVo.XrayProtocol == "vmess") {
		nodeOneVo.Uuid = util.GenerateUUID(*account.Pass)
		nodeOneVo.AlterId = 0
	}
	return nodeOneVo, nil
}

func CreateNode(token string, nodeCreateDto dto.NodeCreateDto) error {
	// 校验端口
	if nodeCreateDto.Port != nil && (*nodeCreateDto.Port <= 100 || *nodeCreateDto.Port >= 30000) {
		return errors.New(constant.PortRangeError)
	}

	// 校验名称
	countName, err := dao.CountNodeByNameAndNodeServerId(nil, nodeCreateDto.Name, nil)
	if err != nil {
		return err
	}
	if countName > 0 {
		return errors.New(constant.NodeNameExist)
	}

	nodeServer, err := dao.SelectNodeServerById(nodeCreateDto.NodeServerId)
	if err != nil {
		return err
	}

	systemName := constant.SystemName
	systemConfig, err := SelectSystemByName(&systemName)
	if err != nil {
		return err
	}

	var nodeId uint
	var mutex sync.Mutex
	defer mutex.Unlock()
	if mutex.TryLock() {
		// Grpc添加节点
		if err = GrpcAddNode(token, *nodeServer.Ip, *nodeServer.GrpcPort, &core.NodeAddDto{
			NodeTypeId: uint64(*nodeCreateDto.NodeTypeId),
			Port:       uint64(*nodeCreateDto.Port),
			Domain:     *nodeCreateDto.Domain,

			//  Xray
			XrayTemplate:       systemConfig.XrayTemplate,
			XrayFlow:           *nodeCreateDto.XrayFlow,
			XraySSMethod:       *nodeCreateDto.XraySSMethod,
			XrayProtocol:       *nodeCreateDto.XrayProtocol,
			XraySettings:       *nodeCreateDto.XraySettings,
			XrayStreamSettings: *nodeCreateDto.XrayStreamSettings,
			XrayTag:            *nodeCreateDto.XrayTag,
			XraySniffing:       *nodeCreateDto.XraySniffing,
			XrayAllocate:       *nodeCreateDto.XrayAllocate,
			// Trojan Go
			TrojanGoSni:             *nodeCreateDto.TrojanGoSni,
			TrojanGoMuxEnable:       uint64(*nodeCreateDto.TrojanGoMuxEnable),
			TrojanGoWebsocketEnable: uint64(*nodeCreateDto.TrojanGoWebsocketEnable),
			TrojanGoWebsocketPath:   *nodeCreateDto.TrojanGoWebsocketPath,
			TrojanGoWebsocketHost:   *nodeCreateDto.TrojanGoWebsocketHost,
			TrojanGoSSEnable:        uint64(*nodeCreateDto.TrojanGoSsEnable),
			TrojanGoSSMethod:        *nodeCreateDto.TrojanGoSsMethod,
			TrojanGoSSPassword:      *nodeCreateDto.TrojanGoSsPassword,
			// Hysteria
			HysteriaProtocol: *nodeCreateDto.HysteriaProtocol,
			HysteriaUpMbps:   int64(*nodeCreateDto.HysteriaUpMbps),
			HysteriaDownMbps: int64(*nodeCreateDto.HysteriaDownMbps),
		}); err != nil {
			go func() {
				for {
					select {
					case <-time.After(8 * time.Second):
						_ = GrpcRemoveNode(token, *nodeServer.Ip, *nodeServer.GrpcPort, *nodeCreateDto.Port, *nodeCreateDto.NodeTypeId)
						return
					}
				}
			}()

			return err
		}
		// 数据插入到数据库中
		if *nodeCreateDto.NodeTypeId == constant.Xray {
			nodeXray := module.NodeXray{
				Protocol:       nodeCreateDto.XrayProtocol,
				XrayFlow:       nodeCreateDto.XrayFlow,
				XraySSMethod:   nodeCreateDto.XraySSMethod,
				Settings:       nodeCreateDto.XraySettings,
				StreamSettings: nodeCreateDto.XrayStreamSettings,
				Tag:            nodeCreateDto.XrayTag,
				Sniffing:       nodeCreateDto.XraySniffing,
				Allocate:       nodeCreateDto.XrayAllocate,
			}
			nodeId, err = dao.CreateNodeXray(&nodeXray)
			if err != nil {
				return err
			}
		} else if *nodeCreateDto.NodeTypeId == constant.TrojanGo {
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
				return err
			}
		} else if *nodeCreateDto.NodeTypeId == constant.Hysteria {
			hysteria := module.NodeHysteria{
				Protocol: nodeCreateDto.HysteriaProtocol,
				UpMbps:   nodeCreateDto.HysteriaUpMbps,
				DownMbps: nodeCreateDto.HysteriaDownMbps,
			}
			nodeId, err = dao.CreateNodeHysteria(&hysteria)
			if err != nil {
				return err
			}
		}

		// 在主表中插入数据
		node := module.Node{
			NodeServerId:       nodeCreateDto.NodeServerId,
			NodeSubId:          &nodeId,
			NodeTypeId:         nodeCreateDto.NodeTypeId,
			Name:               nodeCreateDto.Name,
			NodeServerIp:       nodeServer.Ip,
			NodeServerGrpcPort: nodeServer.GrpcPort,
			Domain:             nodeCreateDto.Domain,
			Port:               nodeCreateDto.Port,
		}
		if err = dao.CreateNode(&node); err != nil {
			return err
		}
	}
	return nil
}

func SelectNodePage(queryName *string, nodeServerId *uint, pageNum *uint, pageSize *uint, c *gin.Context) (*vo.NodePageVo, error) {
	nodePage, total, err := dao.SelectNodePage(queryName, nodeServerId, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	nodeBos := make([]bo.NodeBo, 0)
	for _, item := range *nodePage {
		nodeBo := bo.NodeBo{
			Id:                 *item.Id,
			NodeServerId:       *item.NodeServerId,
			NodeSubId:          *item.NodeSubId,
			NodeTypeId:         *item.NodeTypeId,
			Name:               *item.Name,
			NodeServerIp:       *item.NodeServerIp,
			NodeServerGrpcPort: *item.NodeServerGrpcPort,
			Domain:             *item.Domain,
			Port:               *item.Port,
			CreateTime:         *item.CreateTime,
		}
		nodeBos = append(nodeBos, nodeBo)
	}

	account := util.GetCurrentAccount(c)
	if util.IsAdmin(account.Roles) {
		token := util.GetToken(c)
		splitNodeBos := util.SplitArr(nodeBos, 2)
		var nodeMap sync.Map
		var wg sync.WaitGroup
		for i := range splitNodeBos {
			indexI := i
			wg.Add(1)
			go func() {
				for j := range splitNodeBos[indexI] {
					var ip = splitNodeBos[indexI][j].NodeServerIp
					var grpcPort = splitNodeBos[indexI][j].NodeServerGrpcPort
					status, ok := nodeMap.Load(ip)
					if ok {
						splitNodeBos[indexI][j].Status = status.(int)
					} else {
						var status = 0
						success, err := core.Ping(token, ip, grpcPort)
						if err != nil {
							status = -1
						} else {
							if success {
								status = 1
							}
						}
						splitNodeBos[indexI][j].Status = status
						nodeMap.Store(ip, status)
					}
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}

	nodeVos := make([]vo.NodeVo, 0)
	for _, item := range nodeBos {
		nodeVo := vo.NodeVo{
			Id:           item.Id,
			NodeServerId: item.NodeServerId,
			NodeSubId:    item.NodeSubId,
			NodeTypeId:   item.NodeTypeId,
			Name:         item.Name,
			Domain:       item.Domain,
			Port:         item.Port,
			CreateTime:   item.CreateTime,
			Status:       item.Status,
		}
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

// DeleteNodeById 删除远程节点 删除分表 删除主表
func DeleteNodeById(token string, id *uint) error {
	var mutex sync.Mutex
	defer mutex.TryLock()
	if mutex.TryLock() {
		node, err := dao.SelectNodeById(id)
		if err != nil {
			return err
		}
		_ = GrpcRemoveNode(token, *node.NodeServerIp, *node.NodeServerGrpcPort, *node.Port, *node.NodeTypeId)
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
		if err = dao.DeleteNodeById(id); err != nil {
			return err
		}
	}
	return nil
}

func UpdateNodeById(token string, nodeUpdateDto *dto.NodeUpdateDto) error {
	// 校验端口
	if nodeUpdateDto.Port != nil && (*nodeUpdateDto.Port <= 100 || *nodeUpdateDto.Port >= 30000) {
		return errors.New(constant.PortRangeError)
	}

	// 校验名称
	count, err := dao.CountNodeByNameAndNodeServerId(nodeUpdateDto.Id, nodeUpdateDto.Name, nil)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(constant.NodeNameExist)
	}

	nodeServer, err := dao.SelectNodeServerById(nodeUpdateDto.NodeServerId)
	if err != nil {
		return err
	}

	systemName := constant.SystemName
	systemConfig, err := SelectSystemByName(&systemName)
	if err != nil {
		return err
	}

	var mutex sync.Mutex
	defer mutex.Unlock()
	if mutex.TryLock() {
		nodeEntity, err := dao.SelectNodeById(nodeUpdateDto.Id)
		if err != nil {
			return err
		}
		// Grpc的操作
		if err = GrpcRemoveNode(token, *nodeEntity.NodeServerIp, *nodeEntity.NodeServerGrpcPort, *nodeEntity.Port, *nodeEntity.NodeTypeId); err != nil {
			return err
		}
		if err = GrpcAddNode(token, *nodeServer.Ip, *nodeServer.GrpcPort, &core.NodeAddDto{
			NodeTypeId: uint64(*nodeUpdateDto.NodeTypeId),
			Port:       uint64(*nodeUpdateDto.Port),
			Domain:     *nodeUpdateDto.Domain,

			//  Xray
			XrayTemplate:       systemConfig.XrayTemplate,
			XrayProtocol:       *nodeUpdateDto.XrayProtocol,
			XrayFlow:           *nodeUpdateDto.XrayFlow,
			XraySSMethod:       *nodeUpdateDto.XraySSMethod,
			XraySettings:       *nodeUpdateDto.XraySettings,
			XrayStreamSettings: *nodeUpdateDto.XrayStreamSettings,
			XrayTag:            *nodeUpdateDto.XrayTag,
			XraySniffing:       *nodeUpdateDto.XraySniffing,
			XrayAllocate:       *nodeUpdateDto.XrayAllocate,
			// Trojan Go
			TrojanGoSni:             *nodeUpdateDto.TrojanGoSni,
			TrojanGoMuxEnable:       uint64(*nodeUpdateDto.TrojanGoMuxEnable),
			TrojanGoWebsocketEnable: uint64(*nodeUpdateDto.TrojanGoWebsocketEnable),
			TrojanGoWebsocketPath:   *nodeUpdateDto.TrojanGoWebsocketPath,
			TrojanGoWebsocketHost:   *nodeUpdateDto.TrojanGoWebsocketHost,
			TrojanGoSSEnable:        uint64(*nodeUpdateDto.TrojanGoSsEnable),
			TrojanGoSSMethod:        *nodeUpdateDto.TrojanGoSsMethod,
			TrojanGoSSPassword:      *nodeUpdateDto.TrojanGoSsPassword,
			// Hysteria
			HysteriaProtocol: *nodeUpdateDto.HysteriaProtocol,
			HysteriaUpMbps:   int64(*nodeUpdateDto.HysteriaUpMbps),
			HysteriaDownMbps: int64(*nodeUpdateDto.HysteriaDownMbps),
		}); err != nil {
			_ = GrpcRemoveNode(token, *nodeEntity.NodeServerIp, *nodeEntity.NodeServerGrpcPort, *nodeEntity.Port, *nodeEntity.NodeTypeId)
			return err
		}

		if nodeUpdateDto.NodeTypeId == nodeEntity.NodeTypeId {
			// 没有修改节点类型的情况
			if *nodeEntity.NodeTypeId == constant.Xray {
				nodeXray := module.NodeXray{
					Id:             nodeEntity.NodeSubId,
					Protocol:       nodeUpdateDto.XrayProtocol,
					XrayFlow:       nodeUpdateDto.XrayFlow,
					XraySSMethod:   nodeUpdateDto.XraySSMethod,
					Settings:       nodeUpdateDto.XraySettings,
					StreamSettings: nodeUpdateDto.XrayStreamSettings,
					Tag:            nodeUpdateDto.XrayTag,
					Sniffing:       nodeUpdateDto.XraySniffing,
					Allocate:       nodeUpdateDto.XrayAllocate,
				}
				if err = dao.UpdateNodeXrayById(&nodeXray); err != nil {
					return err
				}
			} else if *nodeEntity.NodeTypeId == constant.TrojanGo {
				nodeTrojanGo := module.NodeTrojanGo{
					Id:              nodeEntity.NodeSubId,
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
			} else if *nodeEntity.NodeTypeId == constant.Hysteria {
				nodeHysteria := module.NodeHysteria{
					Id:       nodeEntity.NodeSubId,
					Protocol: nodeUpdateDto.HysteriaProtocol,
					UpMbps:   nodeUpdateDto.HysteriaUpMbps,
					DownMbps: nodeUpdateDto.HysteriaDownMbps,
				}
				if err = dao.UpdateNodeHysteriaById(&nodeHysteria); err != nil {
					return err
				}
			}
			if *nodeEntity.NodeServerId != *nodeUpdateDto.NodeServerId ||
				*nodeEntity.Name != *nodeUpdateDto.Name ||
				*nodeEntity.NodeServerIp != *nodeServer.Ip ||
				*nodeEntity.Domain != *nodeUpdateDto.Domain ||
				*nodeEntity.Port != *nodeUpdateDto.Port {
				node := module.Node{
					Id:           nodeUpdateDto.Id,
					NodeServerId: nodeUpdateDto.NodeServerId,
					Name:         nodeUpdateDto.Name,
					NodeServerIp: nodeServer.Ip,
					Domain:       nodeUpdateDto.Domain,
					Port:         nodeUpdateDto.Port,
				}
				if err = dao.UpdateNodeById(&node); err != nil {
					return err
				}
			}
		} else {
			// 修改了节点类型的情况 需要删除分库的数据，然后重新再插入
			if *nodeEntity.NodeTypeId == 1 {
				if err = dao.DeleteNodeXrayById(nodeEntity.NodeSubId); err != nil {
					return err
				}
			} else if *nodeEntity.NodeTypeId == 2 {
				if err = dao.DeleteNodeTrojanGoById(nodeEntity.NodeSubId); err != nil {
					return err
				}
			} else if *nodeEntity.NodeTypeId == 3 {
				if err = dao.DeleteNodeHysteriaById(nodeEntity.NodeSubId); err != nil {
					return err
				}
			}

			// 修改了节点类型
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
				Id:                 nodeUpdateDto.Id,
				NodeServerId:       nodeUpdateDto.NodeServerId,
				NodeSubId:          &nodeId,
				NodeTypeId:         nodeUpdateDto.NodeTypeId,
				Name:               nodeUpdateDto.Name,
				NodeServerIp:       nodeServer.Ip,
				NodeServerGrpcPort: nodeServer.GrpcPort,
				Domain:             nodeUpdateDto.Domain,
				Port:               nodeUpdateDto.Port,
			}
			if err = dao.UpdateNodeById(&node); err != nil {
				return err
			}
		}
	}
	return nil
}

func NodeQRCode(accountId *uint, username *string, id *uint) ([]byte, error) {
	nodeUrl, nodeTypeId, err := NodeURL(accountId, username, id)
	if err != nil {
		return nil, err
	}
	if nodeTypeId == constant.NaiveProxy {
		nodeUrl = strings.TrimPrefix(nodeUrl, "naive+https://")
		nodeUrl = fmt.Sprintf("https://%s", base64.StdEncoding.EncodeToString([]byte(nodeUrl)))
	}
	// 生成二维码
	qrCode, err := qrcode.Encode(nodeUrl, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	return qrCode, nil
}

// NodeURL
// xray: https://github.com/XTLS/Xray-core/issues/91
// trojan-go: https://p4gefau1t.github.io/trojan-go/developer/url/
// hysteria:https://github.com/HyNetwork/hysteria/wiki/URI-Scheme
func NodeURL(accountId *uint, username *string, id *uint) (string, uint, error) {

	node, err := dao.SelectNodeById(id)
	if err != nil {
		return "", 0, errors.New(constant.NodeURLError)
	}

	nodeType, err := dao.SelectNodeTypeById(node.NodeTypeId)
	if err != nil {
		return "", 0, errors.New(constant.NodeURLError)
	}

	password, err := dao.SelectConnectPassword(accountId, nil)
	if err != nil {
		return "", 0, errors.New(constant.NodeURLError)
	}

	// 构建URL
	var headBuilder strings.Builder

	if *nodeType.Id == constant.Xray {
		nodeXray, err := dao.SelectNodeXrayById(node.NodeSubId)
		if err != nil {
			return "", 0, errors.New(constant.NodeURLError)
		}
		streamSettings := bo.StreamSettings{}
		if nodeXray.StreamSettings != nil && *nodeXray.StreamSettings != "" {
			if err := json.Unmarshal([]byte(*nodeXray.StreamSettings), &streamSettings); err != nil {
				return "", 0, errors.New(constant.NodeURLError)
			}
		}
		settings := bo.Settings{}
		if nodeXray.Settings != nil && *nodeXray.Settings != "" {
			if err := json.Unmarshal([]byte(*nodeXray.Settings), &settings); err != nil {
				return "", 0, errors.New(constant.NodeURLError)
			}
		}

		connectPass := password

		if *nodeXray.Protocol == "vless" || *nodeXray.Protocol == "vmess" || *nodeXray.Protocol == "trojan" {
			if *nodeXray.Protocol == "vless" || *nodeXray.Protocol == "vmess" {
				connectPass = util.GenerateUUID(password)
			}
			headBuilder.WriteString(fmt.Sprintf("%s://%s@%s:%d?type=%s&security=%s", *nodeXray.Protocol,
				url.PathEscape(connectPass), *node.Domain, *node.Port,
				streamSettings.Network, streamSettings.Security))
			if *nodeXray.Protocol == "vmess" {
				headBuilder.WriteString("&alterId=0")
				if settings.Encryption == "none" {
					headBuilder.WriteString("&encryption=none")
				}
			}

			if *nodeXray.Protocol == "vless" || *nodeXray.Protocol == "trojan" {
				headBuilder.WriteString(fmt.Sprintf("&flow=%s", *nodeXray.XrayFlow))
			}
			if streamSettings.Network == "ws" {
				if streamSettings.WsSettings.Path != "" {
					headBuilder.WriteString(fmt.Sprintf("&path=%s", streamSettings.WsSettings.Path))
				}
				if streamSettings.WsSettings.Host != "" {
					headBuilder.WriteString(fmt.Sprintf("&host=%s", streamSettings.WsSettings.Host))
				}
			}
		} else if *nodeXray.Protocol == "shadowsocks" {
			headBuilder.WriteString(fmt.Sprintf("ss://%s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s@%s:%d", *nodeXray.XraySSMethod,
				connectPass, *node.Domain, *node.Port)))))
		}
	} else if *nodeType.Id == constant.TrojanGo {
		nodeTrojanGo, err := dao.SelectNodeTrojanGoById(node.NodeSubId)
		if err != nil {
			return "", 0, errors.New(constant.NodeURLError)
		}
		headBuilder.WriteString(fmt.Sprintf("trojan-go://%s@%s:%d?", url.PathEscape(password),
			*node.Domain, *node.Port))
		var sni string
		if nodeTrojanGo.Sni != nil && *nodeTrojanGo.Sni != "" {
			sni = *nodeTrojanGo.Sni
		} else {
			sni = *node.Domain
		}
		headBuilder.WriteString(fmt.Sprintf("sni=%s", url.PathEscape(sni)))
		if nodeTrojanGo.WebsocketEnable != nil && *nodeTrojanGo.WebsocketEnable != 0 &&
			nodeTrojanGo.WebsocketPath != nil && *nodeTrojanGo.WebsocketPath != "" {
			headBuilder.WriteString(fmt.Sprintf("&type=%s", url.PathEscape("ws")))
			headBuilder.WriteString(fmt.Sprintf("&path=%s",
				url.PathEscape(fmt.Sprintf("%s", *nodeTrojanGo.WebsocketPath))))
			if nodeTrojanGo.WebsocketHost != nil && *nodeTrojanGo.WebsocketHost != "" {
				headBuilder.WriteString(fmt.Sprintf("&host=%s",
					url.PathEscape(fmt.Sprintf("%s", *nodeTrojanGo.WebsocketHost))))
			}
			if nodeTrojanGo.SsEnable != nil && *nodeTrojanGo.SsEnable != 0 {
				headBuilder.WriteString(fmt.Sprintf("&encryption=%s", url.PathEscape(
					fmt.Sprintf("ss;%s:%s", *nodeTrojanGo.SsMethod, *nodeTrojanGo.SsPassword))))
			}
		}
	} else if *nodeType.Id == constant.Hysteria {
		nodeHysteria, err := dao.SelectNodeHysteriaById(node.NodeSubId)
		if err != nil {
			return "", 0, errors.New(constant.NodeURLError)
		}
		headBuilder.WriteString(fmt.Sprintf("hysteria://%s:%d?protocol=%s&auth=%s&upmbps=%d&downmbps=%d",
			*node.Domain,
			*node.Port,
			*nodeHysteria.Protocol,
			password,
			*nodeHysteria.UpMbps,
			*nodeHysteria.DownMbps))
	} else if *nodeType.Id == constant.NaiveProxy {
		headBuilder.WriteString(fmt.Sprintf("naive+https://%s:%s@%s:%d", *username, password, *node.Domain, *node.Port))
	}

	if node.Name != nil && *node.Name != "" {
		headBuilder.WriteString(fmt.Sprintf("#%s", url.PathEscape(*node.Name)))
	}
	return headBuilder.String(), *nodeType.Id, nil
}

func CountNode() (int, error) {
	return dao.CountNode()
}

func GrpcAddNode(token string, ip string, grpcPort uint, nodeAddDto *core.NodeAddDto) error {
	if err := core.AddNode(token, ip, grpcPort, nodeAddDto); err != nil {
		return err
	}
	return nil
}

func GrpcRemoveNode(token string, ip string, grpcPort uint, port uint, nodeTypeId uint) error {
	if err := core.RemoveNode(token, ip, grpcPort, &core.NodeRemoveDto{
		NodeTypeId: uint64(nodeTypeId),
		Port:       uint64(port),
	}); err != nil {
		return err
	}
	return nil
}
