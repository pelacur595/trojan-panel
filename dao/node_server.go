package dao

import (
	"errors"
	"fmt"
	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/sirupsen/logrus"
	"trojan-panel/module"
	"trojan-panel/module/constant"
)

func SelectNodeServerById(id *uint) (*module.NodeServer, error) {
	var nodeServer module.NodeServer
	where := map[string]interface{}{"id": *id}
	selectFields := []string{"id", "ip", "`name`", "create_time"}
	buildSelect, values, err := builder.BuildSelect("node_server", where, selectFields)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	rows, err := db.Query(buildSelect, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	if err = scanner.Scan(rows, &nodeServer); err == scanner.ErrEmptyResult {
		return nil, errors.New(constant.NodeNotExist)
	} else if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	return &nodeServer, nil
}

func CreateNodeServer(nodeServer *module.NodeServer) error {
	nodeServerEntity := map[string]interface{}{
		"ip":   *nodeServer.Ip,
		"name": *nodeServer.Name,
	}

	var data []map[string]interface{}
	data = append(data, nodeServerEntity)
	buildInsert, values, err := builder.BuildInsert("node_server", data)
	if err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}
	if _, err = db.Exec(buildInsert, values...); err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}
	return nil
}

func SelectNodeServerPage(queryName *string, pageNum *uint, pageSize *uint) (*[]module.NodeServer, uint, error) {
	var (
		total       uint
		nodeServers []module.NodeServer
	)

	// 查询总数
	var whereCount = map[string]interface{}{}
	if queryName != nil && *queryName != "" {
		whereCount["name like"] = fmt.Sprintf("%%%s%%", *queryName)
	}
	selectFieldsCount := []string{"count(1)"}
	buildSelect, values, err := builder.BuildSelect("node_server", whereCount, selectFieldsCount)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, 0, errors.New(constant.SysError)
	}
	if err = db.QueryRow(buildSelect, values...).Scan(&total); err != nil {
		logrus.Errorln(err.Error())
		return nil, 0, errors.New(constant.SysError)
	}

	// 分页查询
	where := map[string]interface{}{
		"_orderby": "create_time desc",
		"_limit":   []uint{(*pageNum - 1) * *pageSize, *pageSize}}
	if queryName != nil && *queryName != "" {
		where["name like"] = fmt.Sprintf("%%%s%%", *queryName)
	}
	selectFields := []string{"id", "`ip`", "name", "create_time"}
	selectSQL, values, err := builder.BuildSelect("node_server", where, selectFields)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, 0, errors.New(constant.SysError)
	}

	rows, err := db.Query(selectSQL, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, 0, errors.New(constant.SysError)
	}
	defer rows.Close()

	if err = scanner.Scan(rows, &nodeServers); err != nil {
		logrus.Errorln(err.Error())
		return nil, 0, errors.New(constant.SysError)
	}
	return &nodeServers, total, nil
}

func DeleteNodeServerById(id *uint) error {
	buildDelete, values, err := builder.BuildDelete("node_server", map[string]interface{}{"id": *id})
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

func UpdateNodeServerById(nodeServer *module.NodeServer) error {
	where := map[string]interface{}{"id": *nodeServer.Id}
	update := map[string]interface{}{}
	if nodeServer.Name != nil {
		update["name"] = *nodeServer.Name
	}
	if nodeServer.Ip != nil {
		update["ip"] = *nodeServer.Ip
	}
	if len(update) > 0 {
		buildUpdate, values, err := builder.BuildUpdate("node_server", where, update)
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

func CountNodeServer() (int, error) {
	return CountNodeServerByName(nil, nil)
}

func CountNodeServerByName(id *uint, queryName *string) (int, error) {
	var count int

	var whereCount = map[string]interface{}{}
	if id != nil {
		whereCount["id <>"] = *id
	}
	if queryName != nil {
		whereCount["name"] = *queryName
	}

	selectFields := []string{"count(1)"}
	buildSelect, values, err := builder.BuildSelect("node_server", whereCount, selectFields)
	if err != nil {
		logrus.Errorln(err.Error())
		return 0, errors.New(constant.SysError)
	}

	if err = db.QueryRow(buildSelect, values...).Scan(&count); err != nil {
		logrus.Errorln(err.Error())
		return 0, errors.New(constant.SysError)
	}
	return count, nil
}

func SelectNodeServerList(ip *string, name *string) ([]module.NodeServer, error) {
	var nodeServers []module.NodeServer
	where := map[string]interface{}{
		"_orderby": "create_time desc"}
	if ip != nil && *ip != "" {
		where["ip like"] = fmt.Sprintf("%%%s%%", *ip)
	}
	if name != nil && *name != "" {
		where["name like"] = fmt.Sprintf("%%%s%%", *name)
	}
	selectFields := []string{"id", "`ip`", "name", "create_time"}
	selectSQL, values, err := builder.BuildSelect("node_server", where, selectFields)
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

	if err = scanner.Scan(rows, &nodeServers); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	return nodeServers, nil
}
