package dao

import (
	"database/sql"
	"github.com/didi/gendry/manager"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
	"trojan/core"
)

var db *sql.DB

// 初始化数据库
func InitDB() {
	mySQLConfig := core.Config.MySQLConfig
	var err error

	db, err = manager.
		New("trojan", mySQLConfig.User, mySQLConfig.Password, mySQLConfig.Host).
		Set(
			manager.SetCharset("utf8"),
			manager.SetAllowCleartextPasswords(true),
			manager.SetInterpolateParams(true),
			manager.SetTimeout(1*time.Second),
			manager.SetReadTimeout(1*time.Second)).
		Port(mySQLConfig.Port).Open(true)

	if err != nil {
		logrus.Errorf("数据库连接异常 err: %v\n", err)
		panic(err)
	}
	//if err = ExecSql(sqlStr); err != nil {
	//	logrus.Errorf("数据库导入失败 err: %v\n", err)
	//	panic(err)
	//}
}

func ExecSql(sqlStr string) error {
	sqls := strings.Split(strings.Replace(sqlStr, "\r\n", "\n", -1), ";\n")
	for _, s := range sqls {
		s = strings.TrimSpace(s)
		if s != "" {
			if _, err := db.Exec(s); err != nil {
				logrus.Errorf("sql执行失败 err: %v\n", err)
				return err
			}
		}
	}
	return nil
}
