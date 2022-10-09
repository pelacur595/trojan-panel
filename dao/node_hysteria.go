package dao

import (
	"errors"
	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/sirupsen/logrus"
	"trojan-panel/module"
	"trojan-panel/module/constant"
)

func SelectNodeHysteriaById(id *uint) (*module.NodeHysteria, error) {
	var nodeHysteria module.NodeHysteria
	where := map[string]interface{}{"id": *id}
	selectFields := []string{"id", "protocol", "up_mbps", "down_mbps"}
	buildSelect, values, err := builder.BuildSelect("node_hysteria", where, selectFields)
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

	if err = scanner.Scan(rows, &nodeHysteria); err == scanner.ErrEmptyResult {
		return nil, errors.New(constant.NodeNotExist)
	} else if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	return &nodeHysteria, nil
}

func CreateNodeHysteria(nodeHysteria *module.NodeHysteria) (uint, error) {
	nodeHysteriaCreate := map[string]interface{}{}
	if nodeHysteria.Protocol != nil && *nodeHysteria.Protocol != "" {
		nodeHysteriaCreate["protocol"] = nodeHysteria.Protocol
	}
	if nodeHysteria.UpMbps != nil {
		nodeHysteriaCreate["up_mbps"] = nodeHysteria.UpMbps
	}
	if nodeHysteria.DownMbps != nil {
		nodeHysteriaCreate["down_mbps"] = nodeHysteria.DownMbps
	}
	if len(nodeHysteriaCreate) > 0 {
		var data []map[string]interface{}
		data = append(data, nodeHysteriaCreate)
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

		if _, err = db.Exec(buildUpdate, values...); err != nil {
			logrus.Errorln(err.Error())
			return errors.New(constant.SysError)
		}
	}
	return nil
}

func DeleteNodeHysteriaById(id *uint) error {
	buildDelete, values, err := builder.BuildDelete("node_hysteria", map[string]interface{}{"id": *id})
	if err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}

	if _, err = db.Exec(buildDelete, values...); err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}
	return nil
}
