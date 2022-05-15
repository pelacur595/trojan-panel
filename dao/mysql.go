package dao

import (
	"database/sql"
	"github.com/didi/gendry/manager"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"net/url"
	"strings"
	"time"
	"trojan/core"
)

var db *sql.DB

// 初始化数据库
func InitMySQL() {
	mySQLConfig := core.Config.MySQLConfig
	var err error

	db, err = manager.
		New("trojan_panel_db", mySQLConfig.User, mySQLConfig.Password, mySQLConfig.Host).
		Set(
			manager.SetCharset("utf8"),
			manager.SetAllowCleartextPasswords(true),
			manager.SetInterpolateParams(true),
			manager.SetTimeout(1*time.Second),
			manager.SetReadTimeout(1*time.Second),
			manager.SetLoc(url.QueryEscape("Asia/Shanghai"))).
		Port(mySQLConfig.Port).Open(true)

	if err != nil {
		logrus.Errorf("数据库连接异常 err: %v\n", err)
		panic(err)
	}

	var count int
	if err = db.QueryRow("SELECT COUNT(1) FROM information_schema.TABLES WHERE table_schema = 'trojan_panel_db' GROUP BY table_schema;").
		Scan(&count); err != nil && err != sql.ErrNoRows {
		logrus.Errorf("查询数据库异常 err: %v\n", err)
		panic(err)
	}
	if count == 0 {
		if err = SqlInit(sqlInitStr); err != nil {
			logrus.Errorf("数据库导入失败 err: %v\n", err)
			panic(err)
		}
	}
}

func SqlInit(sqlStr string) error {
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

var sqlInitStr = "DROP TABLE IF EXISTS `black_list`;\nCREATE TABLE `black_list` (\n  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',\n  `ip` varchar(64) NOT NULL DEFAULT '' COMMENT 'IP地址',\n  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',\n  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',\n  PRIMARY KEY (`id`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='黑名单';\nLOCK TABLES `black_list` WRITE;\nUNLOCK TABLES;\nDROP TABLE IF EXISTS `casbin_rule`;\nCREATE TABLE `casbin_rule` (\n  `p_type` varchar(32) NOT NULL DEFAULT '',\n  `v0` varchar(255) NOT NULL DEFAULT '',\n  `v1` varchar(255) NOT NULL DEFAULT '',\n  `v2` varchar(255) NOT NULL DEFAULT '',\n  `v3` varchar(255) NOT NULL DEFAULT '',\n  `v4` varchar(255) NOT NULL DEFAULT '',\n  `v5` varchar(255) NOT NULL DEFAULT '',\n  KEY `idx_casbin_rule` (`p_type`,`v0`,`v1`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;\nLOCK TABLES `casbin_rule` WRITE;\nINSERT INTO `casbin_rule` VALUES ('p','sysadmin','/api/users/selectUserById','GET','','',''),('p','sysadmin','/api/users/createUser','POST','','',''),('p','sysadmin','/api/users/getUserInfo','GET','','',''),('p','sysadmin','/api/users/selectUserPage','GET','','',''),('p','sysadmin','/api/users/deleteUserById','POST','','',''),('p','sysadmin','/api/users/updateUserProfile','POST','','',''),('p','sysadmin','/api/users/updateUserById','POST','','',''),('p','sysadmin','/api/users/logout','POST','','',''),('p','sysadmin','/api/role/selectRoleList','GET','','',''),('p','sysadmin','/api/node/selectNodeById','GET','','',''),('p','sysadmin','/api/node/createNode','POST','','',''),('p','sysadmin','/api/node/selectNodePage','GET','','',''),('p','sysadmin','/api/node/deleteNodeById','POST','','',''),('p','sysadmin','/api/node/updateNodeById','POST','','',''),('p','sysadmin','/api/node/nodeQRCode','POST','','',''),('p','sysadmin','/api/node/nodeURL','POST','','',''),('p','sysadmin','/api/nodeType/selectNodeTypeList','GET','','',''),('p','sysadmin','/api/dashboard/panelGroup','GET','','',''),('p','sysadmin','/api/dashboard/trafficRank','GET','','',''),('p','sysadmin','/api/system/selectSystemByName','GET','','',''),('p','sysadmin','/api/system/updateSystemById','POST','','',''),('p','sysadmin','/api/system/uploadWebFile','POST','','',''),('p','sysadmin','/api/blackList/selectBlackListPage','GET','','',''),('p','sysadmin','/api/blackList/deleteBlackListByIp','POST','','',''),('p','sysadmin','/api/blackList/createBlackList','POST','','',''),('p','sysadmin','/api/emailRecord/selectEmailRecordPage','GET','','',''),('p','user','/api/users/getUserInfo','GET','','',''),('p','user','/api/users/updateUserProfile','POST','','',''),('p','user','/api/users/logout','POST','','',''),('p','user','/api/node/selectNodePage','GET','','',''),('p','user','/api/node/nodeQRCode','POST','','',''),('p','user','/api/node/nodeURL','POST','','',''),('p','user','/api/nodeType/selectNodeTypeList','GET','','',''),('p','user','/api/dashboard/panelGroup','GET','','',''),('p','user','/api/dashboard/trafficRank','GET','','','');\nUNLOCK TABLES;\nDROP TABLE IF EXISTS `email_record`;\nCREATE TABLE `email_record` (\n  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',\n  `to_email` varchar(64) NOT NULL DEFAULT '' COMMENT '收件人邮箱',\n  `subject` varchar(64) NOT NULL DEFAULT '' COMMENT '主题',\n  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '内容',\n  `state` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '状态 0/未发送 1/发送成功 -1/发送失败',\n  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',\n  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',\n  PRIMARY KEY (`id`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='邮件发送记录';\nLOCK TABLES `email_record` WRITE;\nUNLOCK TABLES;\nDROP TABLE IF EXISTS `node`;\nCREATE TABLE `node` (\n  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',\n  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',\n  `ip` varchar(64) NOT NULL DEFAULT '' COMMENT 'IP地址',\n  `port` int(10) unsigned NOT NULL DEFAULT '443' COMMENT '端口',\n  `type` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '类型 1/trojan-go 2/hysteria',\n  `websocket_enable` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否开启websocket 0/否 1/是',\n  `websocket_path` varchar(64) NOT NULL DEFAULT 'trojan-panel-websocket-path' COMMENT 'websocket路径',\n  `ss_enable` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否开启ss加密 0/否 1/是',\n  `ss_method` varchar(32) NOT NULL DEFAULT 'AES-128-GCM' COMMENT 'ss加密方式',\n  `ss_password` varchar(64) NOT NULL DEFAULT '' COMMENT 'ss密码',\n  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',\n  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',\n  PRIMARY KEY (`id`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='节点';\nLOCK TABLES `node` WRITE;\nUNLOCK TABLES;\nDROP TABLE IF EXISTS `node_type`;\nCREATE TABLE `node_type` (\n  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',\n  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',\n  `prefix` varchar(50) NOT NULL DEFAULT '' COMMENT '节点url前缀',\n  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',\n  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',\n  PRIMARY KEY (`id`)\n) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='节点类型';\nLOCK TABLES `node_type` WRITE;\nINSERT INTO `node_type` VALUES (1,'trojan-go','trojan-go','2022-04-01 00:00:00','2022-04-01 00:00:00'),(2,'hysteria','hysteria','2022-04-01 00:00:00','2022-04-01 00:00:00');\nUNLOCK TABLES;\nDROP TABLE IF EXISTS `role`;\nCREATE TABLE `role` (\n  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',\n  `name` varchar(10) NOT NULL DEFAULT '' COMMENT '名称',\n  `desc` varchar(10) NOT NULL DEFAULT '' COMMENT '描述',\n  `parent_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '父级id',\n  `path` varchar(100) NOT NULL DEFAULT '' COMMENT '路径',\n  `level` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '等级',\n  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',\n  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',\n  PRIMARY KEY (`id`),\n  KEY `role_name_index` (`name`)\n) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COMMENT='角色';\nLOCK TABLES `role` WRITE;\nINSERT INTO `role` VALUES (1,'sysadmin','系统管理员',0,'',1,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(2,'admin','管理员',1,'1-',2,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(3,'user','普通用户',2,'1-2-',3,'2022-04-01 00:00:00','2022-04-01 00:00:00');\nUNLOCK TABLES;\nDROP TABLE IF EXISTS `system`;\nCREATE TABLE `system` (\n  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',\n  `name` varchar(16) NOT NULL DEFAULT '' COMMENT '系统名称',\n  `open_register` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '开放注册 0/否 1/是',\n  `register_quota` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '新默认流量 单位/byte',\n  `register_expire_days` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '新用户默认过期天数 单位/天',\n  `expire_warn_enable` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否开启到期警告 0/否 1/是',\n  `expire_warn_day` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '到期警告 单位/天',\n  `email_enable` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否开启邮箱功能 0/否 1/是',\n  `email_host` varchar(64) NOT NULL DEFAULT '' COMMENT '系统邮箱设置-host',\n  `email_port` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '系统邮箱设置-port',\n  `email_username` varchar(32) NOT NULL DEFAULT '' COMMENT '系统邮箱设置-username',\n  `email_password` varchar(32) NOT NULL DEFAULT '' COMMENT '系统邮箱设置-password',\n  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',\n  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',\n  PRIMARY KEY (`id`)\n) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='系统设置';\nLOCK TABLES `system` WRITE;\nINSERT INTO `system` VALUES (1,'trojan-panel',1,0,0,0,0,0,'',0,'','','2022-04-01 00:00:00','2022-04-01 00:00:00');\nUNLOCK TABLES;\nDROP TABLE IF EXISTS `users`;\nCREATE TABLE `users` (\n  `id` bigint(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',\n  `password` char(56) NOT NULL COMMENT '密码',\n  `quota` bigint(20) NOT NULL DEFAULT '0' COMMENT '配额 单位/byte',\n  `download` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '下载 单位/byte',\n  `upload` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '上传 单位/byte',\n  `username` varchar(64) NOT NULL DEFAULT '' COMMENT '登录用户名',\n  `pass` varchar(64) NOT NULL DEFAULT '' COMMENT '登录密码',\n  `role_id` bigint(20) unsigned NOT NULL DEFAULT '3' COMMENT '角色id 1/系统管理员 3/普通用户',\n  `deleted` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否禁用 0/正常 1/禁用',\n  `expire_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '过期时间',\n  `email` varchar(64) NOT NULL DEFAULT '' COMMENT '邮箱',\n  `ip_limit` tinyint(2) unsigned NOT NULL DEFAULT '3' COMMENT '限制IP设备数',\n  `upload_speed_limit` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '上传限速 单位/byte',\n  `download_speed_limit` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '下载限速 单位/byte',\n  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',\n  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',\n  PRIMARY KEY (`id`),\n  KEY `password` (`password`)\n) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='账户';\nLOCK TABLES `users` WRITE;\nINSERT INTO `users` VALUES (1,'b4fc1369dd766eca295fb495b0938843becbac59fc5cb273b320aaa5',-1,0,0,'sysadmin','MTIzNDU2',1,0,32472115200000,'',3,0,0,'2022-04-01 00:00:00','2022-04-01 00:00:00');\nUNLOCK TABLES;"
