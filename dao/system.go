package dao

import (
	"errors"
	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/sirupsen/logrus"
	"trojan/module"
	"trojan/module/constant"
	"trojan/module/vo"
)

func SelectSystemByName(name *string) (*vo.SystemVo, error) {
	var system module.System
	buildSelect, values, err := builder.NamedQuery(
		"select id,open_register,register_quota,register_expire_days,email_host,email_port,email_username,email_password from `system` where name = {{name}}",
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

	err = scanner.Scan(rows, &system)
	if err == scanner.ErrEmptyResult {
		return nil, errors.New(constant.SystemNotExist)
	} else if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	systemVo := vo.SystemVo{
		Id:                 *system.Id,
		OpenRegister:       *system.OpenRegister,
		RegisterQuota:      *system.RegisterQuota,
		RegisterExpireDays: *system.RegisterExpireDays,
		EmailHost:          *system.EmailHost,
		EmailPort:          *system.EmailPort,
		EmailUsername:      *system.EmailUsername,
		EmailPassword:      *system.EmailPassword,
	}
	return &systemVo, nil
}

func UpdateSystemById(system *module.System) error {
	where := map[string]interface{}{"id": *system.Id}
	update := map[string]interface{}{}
	if system.OpenRegister != nil {
		update["open_register"] = *system.OpenRegister
	}
	if system.RegisterQuota != nil {
		update["register_quota"] = *system.RegisterQuota
	}
	if system.RegisterExpireDays != nil {
		update["register_expire_days"] = *system.RegisterExpireDays
	}
	if system.EmailHost != nil {
		update["email_host"] = *system.EmailHost
	}
	if system.EmailPort != nil {
		update["email_port"] = *system.EmailPort
	}
	if system.EmailUsername != nil {
		update["email_username"] = *system.EmailUsername
	}
	if system.EmailPassword != nil {
		update["email_password"] = *system.EmailPassword
	}
	if len(update) > 0 {
		buildUpdate, values, err := builder.BuildUpdate("`system`", where, update)
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
