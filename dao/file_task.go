package dao

import (
	"errors"
	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/sirupsen/logrus"
	"trojan-panel/module"
	"trojan-panel/module/constant"
	"trojan-panel/module/vo"
)

func SelectFileTaskPage(taskType *uint, pageNum *uint, pageSize *uint) (*vo.FileTaskPageVo, error) {
	var (
		total     uint
		fileTasks []module.FileTask
	)

	// 查询总数
	var whereCount = map[string]interface{}{}
	if taskType != nil && *taskType != 0 {
		whereCount["`type`"] = *taskType
	}
	selectFieldsCount := []string{"count(1)"}
	buildSelect, values, err := builder.BuildSelect("file_task", whereCount, selectFieldsCount)
	if err = db.QueryRow(buildSelect, values...).Scan(&total); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	// 分页查询
	where := map[string]interface{}{
		"_orderby": "role_id,create_time desc",
		"_limit":   []uint{(*pageNum - 1) * *pageSize, *pageSize}}
	if taskType != nil && *taskType != 0 {
		where["`type`"] = *taskType
	}
	selectFields := []string{"id", "name", "`type`", "status", "create_time"}
	selectSQL, values, err := builder.BuildSelect("file_task", where, selectFields)
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

	if err = scanner.Scan(rows, &fileTasks); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	var fileTaskVos = make([]vo.FileTaskVo, 0)
	for _, item := range fileTasks {
		fileTaskVos = append(fileTaskVos, vo.FileTaskVo{
			Id:         *item.Id,
			Name:       *item.Name,
			Type:       *item.Type,
			Status:     *item.Status,
			CreateTime: *item.CreateTime,
		})
	}

	fileTaskPageVo := vo.FileTaskPageVo{
		BaseVoPage: vo.BaseVoPage{
			PageNum:  *pageNum,
			PageSize: *pageSize,
			Total:    total,
		},
		FileTaskVos: fileTaskVos,
	}
	return &fileTaskPageVo, nil
}

func DeleteFileTaskById(id *uint) error {
	buildDelete, values, err := builder.BuildDelete("file_task", map[string]interface{}{"id": *id})
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
