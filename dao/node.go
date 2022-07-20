package dao

import (
	"errors"
	"fmt"
	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/sirupsen/logrus"
	"trojan/module"
	"trojan/module/constant"
	"trojan/module/vo"
)

func SelectNodeById(id *uint) (*vo.NodeVo, error) {
	var node module.Node
	where := map[string]interface{}{"id": *id}
	selectFields := []string{"id", "`name`", "ip", "port", "sni", "type", "websocket_enable",
		"websocket_path", "ss_enable", "ss_method", "ss_password", "hysteria_protocol",
		"hysteria_up_mbps", "hysteria_down_mbps", "create_time"}
	buildSelect, values, err := builder.BuildSelect("node", where, selectFields)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	rows, err := db.Query(buildSelect, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	defer rows.Close()

	err = scanner.Scan(rows, &node)
	if err == scanner.ErrEmptyResult {
		return nil, errors.New(constant.NodeNotExist)
	} else if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	nodeVo := vo.NodeVo{
		Id:               *node.Id,
		Name:             *node.Name,
		Ip:               *node.Ip,
		Port:             *node.Port,
		Sni:              *node.Sni,
		Type:             *node.Type,
		WebsocketEnable:  *node.WebsocketEnable,
		WebsocketPath:    *node.WebsocketPath,
		SsEnable:         *node.SsEnable,
		SsMethod:         *node.SsMethod,
		SsPassword:       *node.SsPassword,
		HysteriaProtocol: *node.HysteriaProtocol,
		HysteriaUpMbps:   *node.HysteriaUpMbps,
		HysteriaDownMbps: *node.HysteriaDownMbps,
		CreateTime:       *node.CreateTime,
	}
	return &nodeVo, nil
}

func CreateNode(node *module.Node) error {
	nodeEntity := map[string]interface{}{
		"name": *node.Name,
		"ip":   *node.Ip,
		"port": *node.Port,
		"type": *node.Type,
	}
	if node.Sni != nil {
		nodeEntity["sni"] = *node.Sni
	}
	if node.WebsocketEnable != nil {
		nodeEntity["websocket_enable"] = *node.WebsocketEnable
	}
	if node.WebsocketPath != nil {
		nodeEntity["websocket_path"] = *node.WebsocketPath
	}
	if node.SsEnable != nil {
		nodeEntity["ss_enable"] = *node.SsEnable
	}
	if node.SsEnable != nil {
		nodeEntity["ss_method"] = *node.SsMethod
	}
	if node.SsPassword != nil {
		nodeEntity["ss_password"] = *node.SsPassword
	}
	if node.HysteriaProtocol != nil {
		nodeEntity["hysteria_protocol"] = *node.HysteriaProtocol
	}
	if node.HysteriaUpMbps != nil {
		nodeEntity["hysteria_up_mbps"] = *node.HysteriaUpMbps
	}
	if node.HysteriaDownMbps != nil {
		nodeEntity["hysteria_down_mbps"] = *node.HysteriaDownMbps
	}

	var data []map[string]interface{}
	data = append(data, nodeEntity)
	buildInsert, values, err := builder.BuildInsert("node", data)
	if err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}
	if _, err := db.Exec(buildInsert, values...); err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}
	return nil
}

func SelectNodePage(queryName *string, pageNum *uint, pageSize *uint) (*vo.NodePageVo, error) {
	var (
		total uint
		nodes []module.Node
	)

	// 查询总数
	var whereCount = map[string]interface{}{}
	if queryName != nil && *queryName != "" {
		whereCount["name like"] = fmt.Sprintf("%%%s%%", *queryName)
	}
	selectFieldsCount := []string{"count(1)"}
	buildSelect, values, err := builder.BuildSelect("node", whereCount, selectFieldsCount)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	if err := db.QueryRow(buildSelect, values...).Scan(&total); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	// 分页查询
	offset := (*pageNum - 1) * *pageSize
	where := map[string]interface{}{
		"_orderby": "create_time desc",
		"_limit":   []uint{offset, *pageSize}}
	if queryName != nil && *queryName != "" {
		where["name like"] = fmt.Sprintf("%%%s%%", *queryName)
	}
	selectFields := []string{"id", "`name`", "ip", "port", "sni", "type", "websocket_enable",
		"websocket_path", "ss_enable", "ss_method", "ss_password", "hysteria_protocol",
		"hysteria_up_mbps", "hysteria_down_mbps", "create_time"}
	selectSQL, values, err := builder.BuildSelect("node", where, selectFields)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	rows, err := db.Query(selectSQL, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	defer rows.Close()

	if err := scanner.Scan(rows, &nodes); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	var nodeVos = make([]vo.NodeVo, 0)
	for _, item := range nodes {
		nodeVos = append(nodeVos, vo.NodeVo{
			Id:               *item.Id,
			Name:             *item.Name,
			Ip:               *item.Ip,
			Port:             *item.Port,
			Sni:              *item.Sni,
			Type:             *item.Type,
			WebsocketEnable:  *item.WebsocketEnable,
			WebsocketPath:    *item.WebsocketPath,
			SsEnable:         *item.SsEnable,
			SsMethod:         *item.SsMethod,
			SsPassword:       *item.SsPassword,
			HysteriaProtocol: *item.HysteriaProtocol,
			HysteriaUpMbps:   *item.HysteriaUpMbps,
			HysteriaDownMbps: *item.HysteriaDownMbps,
			CreateTime:       *item.CreateTime,
		})
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
	buildDelete, values, err := builder.BuildDelete("node", map[string]interface{}{"id": *id})
	if err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}

	if _, err := db.Exec(buildDelete, values...); err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}
	return nil
}

func UpdateNodeById(node *module.Node) error {
	where := map[string]interface{}{"id": *node.Id}
	update := map[string]interface{}{}
	if node.Name != nil {
		update["`name`"] = *node.Name
	}
	if node.Ip != nil {
		update["ip"] = *node.Ip
	}
	if node.Port != nil {
		update["port"] = *node.Port
	}
	if node.Sni != nil {
		update["sni"] = *node.Sni
	}
	if node.Type != nil {
		update["type"] = *node.Type
	}
	if node.WebsocketEnable != nil {
		update["websocket_enable"] = *node.WebsocketEnable
	}
	if node.WebsocketPath != nil && *node.WebsocketPath != "" {
		update["websocket_path"] = *node.WebsocketPath
	}
	if node.SsEnable != nil {
		update["ss_enable"] = *node.SsEnable
	}
	if node.SsMethod != nil && *node.SsMethod != "" {
		update["ss_method"] = *node.SsMethod
	}
	if node.SsPassword != nil && *node.SsPassword != "" {
		update["ss_password"] = *node.SsPassword
	}
	if node.HysteriaProtocol != nil && *node.HysteriaProtocol != "" {
		update["hysteria_protocol"] = *node.HysteriaProtocol
	}
	if node.HysteriaUpMbps != nil {
		update["hysteria_up_mbps"] = *node.HysteriaUpMbps
	}
	if node.HysteriaDownMbps != nil {
		update["hysteria_down_mbps"] = *node.HysteriaDownMbps
	}
	if len(update) > 0 {
		buildUpdate, values, err := builder.BuildUpdate("node", where, update)
		if err != nil {
			logrus.Errorln(err.Error())
			return errors.New(constant.SysError)
		}

		if _, err := db.Exec(buildUpdate, values...); err != nil {
			logrus.Errorln(err.Error())
			return errors.New(constant.SysError)
		}
	}
	return nil
}

func CountNode() (int, error) {
	return CountNodeByName(nil)
}

func CountNodeByName(queryName *string) (int, error) {
	var count int

	var whereCount = map[string]interface{}{}
	if queryName != nil {
		whereCount["name"] = *queryName
	}

	selectFields := []string{"count(1)"}
	buildSelect, values, err := builder.BuildSelect("node", whereCount, selectFields)
	if err != nil {
		logrus.Errorln(err.Error())
		return 0, errors.New(constant.SysError)
	}

	if err := db.QueryRow(buildSelect, values...).Scan(&count); err != nil {
		logrus.Errorln(err.Error())
		return 0, errors.New(constant.SysError)
	}
	return count, nil
}

func SelectNodeIps() ([]string, error) {
	var ips []string

	selectFields := []string{"ip"}
	buildSelect, values, err := builder.BuildSelect("node", nil, selectFields)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	rows, err := db.Query(buildSelect, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	if err := scanner.Scan(rows, &ips); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	return ips, nil
}
