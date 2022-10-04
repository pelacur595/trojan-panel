package dao

import (
	"errors"
	"fmt"
	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/sirupsen/logrus"
	"trojan/module"
	"trojan/module/constant"
)

func SelectNodeById(id *uint) (*module.Node, error) {
	var node module.Node
	where := map[string]interface{}{"id": *id}
	selectFields := []string{"id", "`node_sub_id`", "node_type_id", "name", "ip", "port", "create_time"}
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

	if err = scanner.Scan(rows, &node); err == scanner.ErrEmptyResult {
		return nil, errors.New(constant.NodeNotExist)
	} else if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	return &node, nil
}

func CreateNode(node *module.Node) error {
	nodeEntity := map[string]interface{}{
		"node_sub_id":  *node.NodeSubId,
		"node_type_id": *node.NodeTypeId,
		"name":         *node.Name,
		"ip":           *node.Ip,
	}
	if node.Port != nil && *node.Port != 0 {
		nodeEntity["port"] = *node.Port
	}

	var data []map[string]interface{}
	data = append(data, nodeEntity)
	buildInsert, values, err := builder.BuildInsert("node", data)
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

func SelectNodePage(queryName *string, pageNum *uint, pageSize *uint) (*[]module.Node, uint, error) {
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
	selectFields := []string{"id", "`node_sub_id`", "node_type_id", "name", "ip", "port", "create_time"}
	selectSQL, values, err := builder.BuildSelect("node", where, selectFields)
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

	if err = scanner.Scan(rows, &nodes); err != nil {
		logrus.Errorln(err.Error())
		return nil, 0, errors.New(constant.SysError)
	}
	return &nodes, total, nil
}

func DeleteNodeById(id *uint) error {
	buildDelete, values, err := builder.BuildDelete("node", map[string]interface{}{"id": *id})
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

func UpdateNodeById(node *module.Node) error {
	where := map[string]interface{}{"id": *node.Id}
	update := map[string]interface{}{}
	if node.Name != nil {
		update["name"] = *node.Name
	}
	if node.Ip != nil {
		update["ip"] = *node.Ip
	}
	if node.Port != nil {
		update["port"] = *node.Port
	}
	if len(update) > 0 {
		buildUpdate, values, err := builder.BuildUpdate("node", where, update)
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

	if err = db.QueryRow(buildSelect, values...).Scan(&count); err != nil {
		logrus.Errorln(err.Error())
		return 0, errors.New(constant.SysError)
	}
	return count, nil
}

func SelectNodesIpAndPort() ([]module.Node, error) {
	var nodes []module.Node

	buildSelect, values, err := builder.BuildSelect("node", nil, []string{"id", "ip", "port"})
	if err != nil {
		logrus.Errorln(err.Error())
		return nodes, errors.New(constant.SysError)
	}
	rows, err := db.Query(buildSelect, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return nodes, errors.New(constant.SysError)
	}
	defer rows.Close()

	if err = scanner.Scan(rows, &nodes); err != nil {
		logrus.Errorln(err.Error())
		return nodes, errors.New(constant.SysError)
	}
	return nodes, nil
}

func SelectNodesIpDistinct() ([]string, error) {
	var ips []string

	buildSelect, values, err := builder.NamedQuery("select distinct ip from node", nil)
	if err != nil {
		logrus.Errorln(err.Error())
		return ips, errors.New(constant.SysError)
	}
	rows, err := db.Query(buildSelect, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return ips, errors.New(constant.SysError)
	}
	defer rows.Close()

	result, err := scanner.ScanMap(rows)
	if err != nil {
		logrus.Errorln(err.Error())
		return ips, errors.New(constant.SysError)
	}
	for _, record := range result {
		ips = append(ips, fmt.Sprintf("%s", record["ip"]))
	}
	return ips, nil
}
