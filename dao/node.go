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
	selectFields := []string{"id", "`name`", "ip", "port", "type", "create_time"}
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
		Id:         *node.Id,
		Name:       *node.Name,
		Ip:         *node.Ip,
		Port:       *node.Port,
		Type:       *node.Type,
		CreateTime: *node.CreateTime,
	}
	return &nodeVo, nil
}

func CreateNode(node *module.Node) error {
	var data []map[string]interface{}
	data = append(data, map[string]interface{}{
		"name": *node.Name,
		"ip":   *node.Ip,
		"port": *node.Port,
		"type": *node.Type,
	})

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
	selectFields := []string{"id", "`name`", "ip", "port",
		"type", "create_time"}
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

	if err = scanner.Scan(rows, &nodes); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	var nodeVos []vo.NodeVo
	for _, item := range nodes {
		nodeVos = append(nodeVos, vo.NodeVo{
			Id:         *item.Id,
			Name:       *item.Name,
			Ip:         *item.Ip,
			Port:       *item.Port,
			Type:       *item.Type,
			CreateTime: *item.CreateTime,
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
	if node.Type != nil {
		update["type"] = *node.Type
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

	selectFields := []string{"count(1) count"}
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
