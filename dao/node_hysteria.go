package dao

import (
	"errors"
	"github.com/didi/gendry/builder"
	"github.com/sirupsen/logrus"
	"trojan/module"
	"trojan/module/constant"
)

func SelectHysteriaById(id *uint) (*module.NodeHysteria, error) {
	var nodeHysteria module.NodeHysteria
	where := map[string]interface{}{"id": *id}
	selectFields := []string{"id", "protocol", "up_mbps", "down_mbps"}
	buildSelect, values, err := builder.BuildSelect("node_trojan_go", where, selectFields)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	if err = db.QueryRow(buildSelect, values...).Scan(&nodeHysteria); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	return &nodeHysteria, nil
}

func CreateNodeHysteria(nodeHysteria *module.NodeHysteria) (uint, error) {
	nodeHysteriaEntity := map[string]interface{}{}
	if nodeHysteria.Protocol != nil && *nodeHysteria.Protocol != "" {
		nodeHysteriaEntity["protocol"] = nodeHysteria.Protocol
	}
	if nodeHysteria.UpMbps != nil {
		nodeHysteriaEntity["up_mbps"] = nodeHysteria.UpMbps
	}
	if nodeHysteria.DownMbps != nil {
		nodeHysteriaEntity["down_mbps"] = nodeHysteria.DownMbps
	}
	if len(nodeHysteriaEntity) > 0 {
		var data []map[string]interface{}
		data = append(data, nodeHysteriaEntity)
		buildInsert, values, err := builder.BuildInsert("node_hysteria", data)
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

func UpdateNodeHysteriaById(nodeHysteria *module.NodeHysteria) error {
	where := map[string]interface{}{"id": *nodeHysteria.Id}
	update := map[string]interface{}{}
	if nodeHysteria.Protocol != nil && *nodeHysteria.Protocol != "" {
		update["protocol"] = *nodeHysteria.Protocol
	}
	if nodeHysteria.UpMbps != nil {
		update["up_mbps"] = *nodeHysteria.UpMbps
	}
	if nodeHysteria.DownMbps != nil {
		update["down_mbps"] = *nodeHysteria.DownMbps
	}
	if len(update) > 0 {
		buildUpdate, values, err := builder.BuildUpdate("node_hysteria", where, update)
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
