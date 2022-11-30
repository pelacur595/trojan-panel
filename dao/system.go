package dao

import (
	"errors"
	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/sirupsen/logrus"
	"trojan-panel/module"
	"trojan-panel/module/constant"
)

func SelectSystemByName(name *string) (*module.System, error) {
	var system module.System
	buildSelect, values, err := builder.NamedQuery(
		"select id,register_config,email_config from `system` where name = {{name}}",
		map[string]interface{}{"name": *name})
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

	if err = scanner.Scan(rows, &system); err == scanner.ErrEmptyResult {
		return nil, errors.New(constant.SystemNotExist)
	} else if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	return &system, nil
}

func UpdateSystemById(system *module.System) error {
	where := map[string]interface{}{"id": *system.Id}
	update := map[string]interface{}{}
	if system.RegisterConfig != nil {
		update["register_config"] = *system.RegisterConfig
	}
	if system.EmailConfig != nil {
		update["email_config"] = *system.EmailConfig
	}
	if len(update) > 0 {
		buildUpdate, values, err := builder.BuildUpdate("`system`", where, update)
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
