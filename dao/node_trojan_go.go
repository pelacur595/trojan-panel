package dao

import (
	"errors"
	"github.com/didi/gendry/builder"
	"github.com/sirupsen/logrus"
	"trojan/module"
	"trojan/module/constant"
)

func SelectNodeTrojanGoById(id *uint) (*module.NodeTrojanGo, error) {
	var nodeTrojanGo module.NodeTrojanGo
	where := map[string]interface{}{"id": *id}
	selectFields := []string{"id", "`sni`", "mux_enable", "websocket_enable", "websocket_path", "websocket_host", "ss_enable", "ss_method", "ss_password"}
	buildSelect, values, err := builder.BuildSelect("node_trojan_go", where, selectFields)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	if err = db.QueryRow(buildSelect, values...).Scan(&nodeTrojanGo); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	return &nodeTrojanGo, nil
}

func CreateNodeTrojanGo(nodeTrojanGo *module.NodeTrojanGo) (uint, error) {
	nodeTrojanGoCreate := map[string]interface{}{}
	if nodeTrojanGo.Sni != nil && *nodeTrojanGo.Sni != "" {
		nodeTrojanGoCreate["sni"] = nodeTrojanGo.Sni
	}
	if nodeTrojanGo.MuxEnable != nil {
		nodeTrojanGoCreate["mux_enable"] = nodeTrojanGo.MuxEnable
	}
	if nodeTrojanGo.WebsocketEnable != nil {
		nodeTrojanGoCreate["websocket_enable"] = nodeTrojanGo.WebsocketEnable
	}
	if nodeTrojanGo.WebsocketPath != nil && *nodeTrojanGo.WebsocketPath != "" {
		nodeTrojanGoCreate["websocket_path"] = nodeTrojanGo.WebsocketPath
	}
	if nodeTrojanGo.SsEnable != nil {
		nodeTrojanGoCreate["ss_enable"] = nodeTrojanGo.SsEnable
	}
	if nodeTrojanGo.SsMethod != nil && *nodeTrojanGo.SsMethod != "" {
		nodeTrojanGoCreate["ss_method"] = nodeTrojanGo.SsMethod
	}
	if nodeTrojanGo.SsPassword != nil && *nodeTrojanGo.SsPassword != "" {
		nodeTrojanGoCreate["ss_password"] = nodeTrojanGo.SsPassword
	}
	if len(nodeTrojanGoCreate) > 0 {
		var data []map[string]interface{}
		data = append(data, nodeTrojanGoCreate)
		buildInsert, values, err := builder.BuildInsert("node_trojan_go", data)
		if err != nil {
			logrus.Errorln(err.Error())
			return 0, errors.New(constant.SysError)
		}
		result, err := db.Exec(buildInsert, values...)
		if err != nil {
			logrus.Errorln(err.Error())
			return 0, errors.New(constant.SysError)
		}
		id, err := result.LastInsertId()
		if err != nil {
			logrus.Errorln(err.Error())
			return 0, errors.New(constant.SysError)
		}
		return uint(id), nil
	}
	return 0, errors.New(constant.SysError)
}

func UpdateNodeTrojanGoById(nodeTrojanGo *module.NodeTrojanGo) error {
	where := map[string]interface{}{"id": *nodeTrojanGo.Id}
	update := map[string]interface{}{}
	if nodeTrojanGo.Sni != nil && *nodeTrojanGo.Sni != "" {
		update["sni"] = *nodeTrojanGo.Sni
	}
	if nodeTrojanGo.MuxEnable != nil {
		update["mux_enable"] = *nodeTrojanGo.MuxEnable
	}
	if nodeTrojanGo.WebsocketEnable != nil {
		update["websocket_enable"] = *nodeTrojanGo.WebsocketEnable
	}
	if nodeTrojanGo.WebsocketPath != nil && *nodeTrojanGo.WebsocketPath != "" {
		update["websocket_path"] = *nodeTrojanGo.WebsocketPath
	}
	if nodeTrojanGo.WebsocketHost != nil && *nodeTrojanGo.WebsocketHost != "" {
		update["websocket_host"] = *nodeTrojanGo.WebsocketHost
	}
	if nodeTrojanGo.SsEnable != nil {
		update["ss_enable"] = *nodeTrojanGo.SsEnable
	}
	if nodeTrojanGo.SsMethod != nil && *nodeTrojanGo.SsMethod != "" {
		update["ss_method"] = *nodeTrojanGo.SsMethod
	}
	if nodeTrojanGo.SsPassword != nil && *nodeTrojanGo.SsPassword != "" {
		update["ss_password"] = *nodeTrojanGo.SsPassword
	}
	if len(update) > 0 {
		buildUpdate, values, err := builder.BuildUpdate("node_trojan_go", where, update)
		if err != nil {
			logrus.Errorln(err.Error())
			return errors.New(constant.SysError)
		}

		if _, err = db.Exec(buildUpdate, values...); err != nil {
			logrus.Errorln(err.Error())
			return errors.New(constant.SysError)
		}
	}
	return nil
}
