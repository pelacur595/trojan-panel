package service

import (
	"errors"
	"fmt"
	redisgo "github.com/gomodule/redigo/redis"
	"github.com/skip2/go-qrcode"
	"net/url"
	"strings"
	"trojan/dao"
	"trojan/dao/redis"
	"trojan/module"
	"trojan/module/constant"
	"trojan/module/dto"
	"trojan/module/vo"
	"trojan/util"
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
	redis.Client.Key.Del("trojan-panel:nodeIps")
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
	if err := dao.DeleteNodeById(id); err != nil {
		return err
	}
	redis.Client.Key.Del("trojan-panel:nodeIps")
	return nil
}

func UpdateNodeById(node *module.Node) error {
	if err := dao.UpdateNodeById(node); err != nil {
		return err
	}
	redis.Client.Key.Del("trojan-panel:nodeIps")
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

	// 构建URL
	var headBuilder strings.Builder
	if nodeTypeVo.Prefix == constant.TrojanGoPrefix {
		headBuilder.WriteString(fmt.Sprintf("%s://%s@%s:%d?", nodeTypeVo.Prefix, url.PathEscape(password),
			*nodeQRCodeDto.Ip, *nodeQRCodeDto.Port))
		if nodeQRCodeDto.Sni != nil && *nodeQRCodeDto.Sni != "" {
			headBuilder.WriteString(fmt.Sprintf("sni=%s", *nodeQRCodeDto.Sni))
		}
		if nodeQRCodeDto.WebsocketEnable != nil && *nodeQRCodeDto.WebsocketEnable != 0 &&
			nodeQRCodeDto.WebsocketPath != nil && *nodeQRCodeDto.WebsocketPath != "" {
			headBuilder.WriteString(fmt.Sprintf("&type=%s", url.PathEscape("ws")))
			headBuilder.WriteString(fmt.Sprintf("&path=%s",
				url.PathEscape(fmt.Sprintf("/%s", *nodeQRCodeDto.WebsocketPath))))
			if nodeQRCodeDto.SsEnable != nil && *nodeQRCodeDto.SsEnable != 0 ||
				nodeQRCodeDto.SsMethod != nil && *nodeQRCodeDto.SsMethod != "" ||
				nodeQRCodeDto.SsPassword != nil && *nodeQRCodeDto.SsPassword != "" {
				headBuilder.WriteString(fmt.Sprintf("&encryption=%s", url.PathEscape(
					fmt.Sprintf("ss;%s:%s", *nodeQRCodeDto.SsMethod, *nodeQRCodeDto.SsPassword))))
			}
		}
	}

	if nodeTypeVo.Prefix == constant.HysteriaPrefix {
		headBuilder.WriteString(fmt.Sprintf("%s://%s:%d?protocol=%s&auth=%s&upmbps=%d&downmbps=%d",
			nodeTypeVo.Prefix,
			*nodeQRCodeDto.Ip,
			*nodeQRCodeDto.Port,
			*nodeQRCodeDto.HysteriaProtocol,
			password,
			*nodeQRCodeDto.HysteriaUpMbps,
			*nodeQRCodeDto.HysteriaDownMbps))
	}
	if nodeQRCodeDto.Name != nil && *nodeQRCodeDto.Name != "" {
		headBuilder.WriteString(fmt.Sprintf("#%s", url.PathEscape(*nodeQRCodeDto.Name)))
	}
	return headBuilder.String(), nil
}

func CountNode() (int, error) {
	nodeCount, err := dao.CountNode()
	if err != nil {
		return 0, err
	}
	return nodeCount, nil
}

func SelectNodeIps() ([]string, error) {
	ips, err := redis.Client.Set.SMembers("trojan-panel:nodeIps").Strings()
	if err != nil && err != redisgo.ErrNil {
		return nil, errors.New(constant.SysError)
	}
	if len(ips) > 0 {
		return ips, nil
	} else {
		nodeIps, err := dao.SelectNodeIps()
		if err != nil {
			return nil, err
		}
		var redisIps []interface{}
		for _, ip := range nodeIps {
			redisIps = append(redisIps, ip)
		}
		redis.Client.Set.SAdd("trojan-panel:nodeIps", redisIps)
		return nodeIps, nil
	}
}
