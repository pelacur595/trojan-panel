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
		New("trojan_panel_db", mySQLConfig.User, mySQLConfig.Password, mySQLConfig.Host).
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

var sqlInitStr = "DROP TABLE IF EXISTS `casbin_rule`;\nCREATE TABLE `casbin_rule` (`p_type` varchar(32) NOT NULL DEFAULT '',`v0` varchar(255) NOT NULL DEFAULT '',`v1` varchar(255) NOT NULL DEFAULT '',`v2` varchar(255) NOT NULL DEFAULT '',`v3` varchar(255) NOT NULL DEFAULT '',`v4` varchar(255) NOT NULL DEFAULT '',`v5` varchar(255) NOT NULL DEFAULT '',KEY `idx_casbin_rule` (`p_type`,`v0`,`v1`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;\nLOCK TABLES `casbin_rule` WRITE;\nINSERT INTO `casbin_rule` VALUES ('p','sysadmin','/api/users/selectUserById','GET','','',''),('p','sysadmin','/api/users/createUser','POST','','',''),('p','sysadmin','/api/users/getUserInfo','GET','','',''),('p','sysadmin','/api/users/selectUserPage','GET','','',''),('p','sysadmin','/api/users/deleteUserById','POST','','',''),('p','sysadmin','/api/users/updateUserPassByUsername','POST','','',''),('p','sysadmin','/api/users/updateUserById','POST','','',''),('p','sysadmin','/api/role/selectRoleList','GET','','',''),('p','sysadmin','/api/node/selectNodeById','GET','','',''),('p','sysadmin','/api/node/createNode','POST','','',''),('p','sysadmin','/api/node/selectNodePage','GET','','',''),('p','sysadmin','/api/node/deleteNodeById','POST','','',''),('p','sysadmin','/api/node/updateNodeById','POST','','',''),('p','sysadmin','/api/node/nodeQRCode','GET','','',''),('p','sysadmin','/api/node/nodeURL','GET','','',''),('p','sysadmin','/api/nodeType/selectNodeTypeList','GET','','',''),('p','sysadmin','/api/trojan-gfw/status','GET','','',''),('p','sysadmin','/api/trojan-gfw/restart','POST','','',''),('p','sysadmin','/api/trojan-gfw/stop','POST','','',''),('p','sysadmin','/api/trojan-go/status','GET','','',''),('p','sysadmin','/api/trojan-go/restart','POST','','',''),('p','sysadmin','/api/trojan-go/stop','POST','','',''),('p','sysadmin','/api/dashboard/panelGroup','GET','','',''),('p','sysadmin','/api/system/selectSystemByName','GET','','',''),('p','sysadmin','/api/system/updateSystemById','POST','','',''),('p','sysadmin','/api/system/uploadWebFile','POST','','',''),('p','user','/api/users/getUserInfo','GET','','',''),('p','user','/api/users/updateUserPassByUsername','POST','','',''),('p','user','/api/node/selectNodePage','GET','','',''),('p','user','/api/node/nodeQRCode','GET','','',''),('p','user','/api/node/nodeURL','GET','','',''),('p','user','/api/nodeType/selectNodeTypeList','GET','','',''),('p','user','/api/dashboard/panelGroup','GET','','','');\nUNLOCK TABLES;\nDROP TABLE IF EXISTS `menu_list`;\nCREATE TABLE `menu_list` (`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',`name` varchar(10) NOT NULL DEFAULT '' COMMENT '名称',`icon` varchar(20) NOT NULL DEFAULT '' COMMENT '图标',`route` varchar(50) NOT NULL DEFAULT '' COMMENT '路由',`order` bigint(20) unsigned NOT NULL DEFAULT '100' COMMENT '排序 越小越靠前',`parent_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '父级id',`path` varchar(100) NOT NULL DEFAULT '' COMMENT '路径',`level` int(11) unsigned NOT NULL DEFAULT '1' COMMENT '等级',`create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',`update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',PRIMARY KEY (`id`),KEY `menu_list_name_index` (`name`)) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COMMENT='菜单';\nLOCK TABLES `menu_list` WRITE;\nINSERT INTO `menu_list` VALUES (1,'仪表板','dashboard','/dashboard',100,0,'',1,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(2,'用户管理','user','/users-manage',100,0,'',1,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(3,'节点管理','node','/node-manage',100,0,'',1,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(4,'用户列表','','/users-manage/user-list',100,2,'2-',2,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(5,'节点列表','','/node-manage/node-list',100,3,'3-',2,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(6,'系统设置','system','/system',100,0,'',1,'2022-04-01 00:00:00','2022-04-01 00:00:00');\nUNLOCK TABLES;\nDROP TABLE IF EXISTS `node`;\nCREATE TABLE `node` (`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',`name` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',`ip` varchar(64) NOT NULL DEFAULT '' COMMENT 'IP地址',`port` int(10) unsigned NOT NULL DEFAULT '443' COMMENT '端口',`type` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '类型 1/trojan-go 2/trojan-gfw',`create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',`update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='节点';\nLOCK TABLES `node` WRITE;\nUNLOCK TABLES;\nDROP TABLE IF EXISTS `node_type`;\nCREATE TABLE `node_type` (`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',`name` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',`prefix` varchar(50) NOT NULL DEFAULT '' COMMENT '节点url前缀',`create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',`update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='节点类型';\nLOCK TABLES `node_type` WRITE;\nINSERT INTO `node_type` VALUES (1,'trojan-go','trojan-go','2022-04-01 00:00:00','2022-04-01 00:00:00'),(2,'trojan-gfw','trojan','2022-04-01 00:00:00','2022-04-01 00:00:00');\nUNLOCK TABLES;\nDROP TABLE IF EXISTS `role`;\nCREATE TABLE `role` (`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',`name` varchar(10) NOT NULL DEFAULT '' COMMENT '名称',`desc` varchar(10) NOT NULL DEFAULT '' COMMENT '描述',`parent_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '父级id',`path` varchar(100) NOT NULL DEFAULT '' COMMENT '路径',`level` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '等级',`create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',`update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',PRIMARY KEY (`id`),KEY `role_name_index` (`name`)) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COMMENT='角色';\nLOCK TABLES `role` WRITE;\nINSERT INTO `role` VALUES (1,'sysadmin','系统管理员',0,'',1,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(2,'admin','管理员',1,'1-',2,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(3,'user','普通用户',2,'1-2-',3,'2022-04-01 00:00:00','2022-04-01 00:00:00');\nUNLOCK TABLES;\nDROP TABLE IF EXISTS `role_menu_list`;\nCREATE TABLE `role_menu_list` (`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',`role_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '角色id',`menu_list_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '菜单id',`create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',`update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COMMENT='角色和菜单关系';\nLOCK TABLES `role_menu_list` WRITE;\nINSERT INTO `role_menu_list` VALUES (1,1,1,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(2,1,2,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(3,1,3,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(4,1,4,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(5,1,5,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(6,1,6,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(7,3,1,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(8,3,3,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(9,3,5,'2022-04-01 00:00:00','2022-04-01 00:00:00');\nUNLOCK TABLES;\nDROP TABLE IF EXISTS `system`;\nCREATE TABLE `system` (`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',`name` varchar(20) NOT NULL DEFAULT '' COMMENT '系统名称',`open_register` tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '开放注册 0/否 1/是',`register_quota` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '新默认流量 单位/byte',`register_expire_days` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '新用户默认过期天数 单位/天',`create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',`update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='系统设置';\nLOCK TABLES `system` WRITE;\nINSERT INTO `system` VALUES (2,'trojan-panel',1,0,0,'2022-04-01 00:00:00','2022-04-01 00:00:00');\nUNLOCK TABLES;\nDROP TABLE IF EXISTS `users`;\nCREATE TABLE `users` (`id` bigint(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',`password` char(56) NOT NULL COMMENT '密码',`quota` bigint(20) NOT NULL DEFAULT '0' COMMENT '配额',`download` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '下载',`upload` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '上传',`username` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',`pass` varchar(64) NOT NULL DEFAULT '' COMMENT '用户密码',`role_id` bigint(20) unsigned NOT NULL DEFAULT '3' COMMENT '角色id 1/系统管理员 3/普通用户',`deleted` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '是否禁用 0/正常 1/禁用',`expire_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '过期时间',`create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',`update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',PRIMARY KEY (`id`),KEY `password` (`password`)) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='账户';\nLOCK TABLES `users` WRITE;\nINSERT INTO `users` VALUES (1,'b4fc1369dd766eca295fb495b0938843becbac59fc5cb273b320aaa5',-1,0,0,'sysadmin','MTIzNDU2',1,0,32472115200000,'2022-04-01 00:00:00','2022-04-01 00:00:00');\nUNLOCK TABLES;\n"
