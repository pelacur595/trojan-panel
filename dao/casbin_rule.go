package dao

import (
	"errors"
	"trojan-panel/module/constant"

	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/sirupsen/logrus"
)

func Casbin() (*casbin.Enforcer, error) {
	a, err := sqladapter.NewAdapter(db, "mysql", "casbin_rule")
	if err != nil {
		logrus.Errorf("casbin初始化失败 err: %v\n", err)
		return nil, errors.New(constant.SysError)
	}
	// 读取conf配置文件
	e, err := casbin.NewEnforcer(constant.RbacModelFilePath, a)
	if err != nil {
		logrus.Errorf("未找到配置文件 error: %v\n", err)
		return nil, errors.New(constant.SysError)
	}
	// 加载规则
	if err := e.LoadPolicy(); err != nil {
		logrus.Errorf("加载规则失败 error: %v\n", err)
		return nil, errors.New(constant.SysError)
	}
	return e, nil
}
