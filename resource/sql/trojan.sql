-- MySQL dump 10.13  Distrib 5.7.35, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: trojan
-- ------------------------------------------------------
-- Server version	5.7.35

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `casbin_rule`
--

DROP TABLE IF EXISTS `casbin_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `casbin_rule` (
  `p_type` varchar(32) NOT NULL DEFAULT '',
  `v0` varchar(255) NOT NULL DEFAULT '',
  `v1` varchar(255) NOT NULL DEFAULT '',
  `v2` varchar(255) NOT NULL DEFAULT '',
  `v3` varchar(255) NOT NULL DEFAULT '',
  `v4` varchar(255) NOT NULL DEFAULT '',
  `v5` varchar(255) NOT NULL DEFAULT '',
  KEY `idx_casbin_rule` (`p_type`,`v0`,`v1`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `casbin_rule`
--

LOCK TABLES `casbin_rule` WRITE;
/*!40000 ALTER TABLE `casbin_rule` DISABLE KEYS */;
INSERT INTO `casbin_rule` VALUES ('p','sysadmin','/api/users/selectUserById','GET','','',''),('p','sysadmin','/api/users/createUser','POST','','',''),('p','sysadmin','/api/users/getUserInfo','GET','','',''),('p','sysadmin','/api/users/selectUserPage','GET','','',''),('p','sysadmin','/api/users/deleteUserById','POST','','',''),('p','sysadmin','/api/users/updateUserPassByUsername','POST','','',''),('p','sysadmin','/api/users/updateUserById','POST','','',''),('p','sysadmin','/api/role/selectRoleList','GET','','',''),('p','sysadmin','/api/node/selectNodeById','GET','','',''),('p','sysadmin','/api/node/createNode','POST','','',''),('p','sysadmin','/api/node/selectNodePage','GET','','',''),('p','sysadmin','/api/node/deleteNodeById','POST','','',''),('p','sysadmin','/api/node/updateNodeById','POST','','',''),('p','sysadmin','/api/node/nodeQRCode','GET','','',''),('p','sysadmin','/api/node/nodeURL','GET','','',''),('p','sysadmin','/api/nodeType/selectNodeTypeList','GET','','',''),('p','sysadmin','/api/trojan-gfw/status','GET','','',''),('p','sysadmin','/api/trojan-gfw/restart','POST','','',''),('p','sysadmin','/api/trojan-gfw/stop','POST','','',''),('p','sysadmin','/api/trojan-go/status','GET','','',''),('p','sysadmin','/api/trojan-go/restart','POST','','',''),('p','sysadmin','/api/trojan-go/stop','POST','','',''),('p','sysadmin','/api/dashboard/panelGroup','GET','','',''),('p','sysadmin','/api/system/selectSystemByName','GET','','',''),('p','sysadmin','/api/system/updateSystemById','POST','','',''),('p','sysadmin','/api/system/uploadWebFile','POST','','',''),('p','user','/api/users/getUserInfo','GET','','',''),('p','user','/api/users/updateUserPassByUsername','POST','','',''),('p','user','/api/node/selectNodePage','GET','','',''),('p','user','/api/node/nodeQRCode','GET','','',''),('p','user','/api/node/nodeURL','GET','','',''),('p','user','/api/nodeType/selectNodeTypeList','GET','','',''),('p','user','/api/dashboard/panelGroup','GET','','','');
/*!40000 ALTER TABLE `casbin_rule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `menu_list`
--

DROP TABLE IF EXISTS `menu_list`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `menu_list` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `name` varchar(10) NOT NULL DEFAULT '' COMMENT '名称',
  `icon` varchar(20) NOT NULL DEFAULT '' COMMENT '图标',
  `route` varchar(50) NOT NULL DEFAULT '' COMMENT '路由',
  `order` bigint(20) unsigned NOT NULL DEFAULT '100' COMMENT '排序 越小越靠前',
  `parent_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '父级id',
  `path` varchar(100) NOT NULL DEFAULT '' COMMENT '路径',
  `level` int(11) unsigned NOT NULL DEFAULT '1' COMMENT '等级',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `menu_list_name_index` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COMMENT='菜单';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `menu_list`
--

LOCK TABLES `menu_list` WRITE;
/*!40000 ALTER TABLE `menu_list` DISABLE KEYS */;
INSERT INTO `menu_list` VALUES (1,'仪表板','dashboard','/dashboard',100,0,'',1,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(2,'用户管理','user','/users-manage',100,0,'',1,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(3,'节点管理','node','/node-manage',100,0,'',1,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(4,'用户列表','','/users-manage/user-list',100,2,'2-',2,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(5,'节点列表','','/node-manage/node-list',100,3,'3-',2,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(6,'系统设置','system','/system',100,0,'',1,'2022-04-01 00:00:00','2022-04-01 00:00:00');
/*!40000 ALTER TABLE `menu_list` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `node`
--

DROP TABLE IF EXISTS `node`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `node` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
  `ip` varchar(128) NOT NULL DEFAULT '' COMMENT 'IP地址',
  `port` int(10) unsigned NOT NULL DEFAULT '443' COMMENT '端口',
  `type` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '类型 1/trojan-go 2/trojan-gfw',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='节点';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `node`
--

LOCK TABLES `node` WRITE;
/*!40000 ALTER TABLE `node` DISABLE KEYS */;
/*!40000 ALTER TABLE `node` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `node_type`
--

DROP TABLE IF EXISTS `node_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `node_type` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
  `prefix` varchar(50) NOT NULL DEFAULT '' COMMENT '节点url前缀',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='节点类型';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `node_type`
--

LOCK TABLES `node_type` WRITE;
/*!40000 ALTER TABLE `node_type` DISABLE KEYS */;
INSERT INTO `node_type` VALUES (1,'trojan-go','trojan-go','2022-04-01 00:00:00','2022-04-01 00:00:00'),(2,'trojan-gfw','trojan','2022-04-01 00:00:00','2022-04-01 00:00:00');
/*!40000 ALTER TABLE `node_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role`
--

DROP TABLE IF EXISTS `role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `role` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `name` varchar(10) NOT NULL DEFAULT '' COMMENT '名称',
  `desc` varchar(10) NOT NULL DEFAULT '' COMMENT '描述',
  `parent_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '父级id',
  `path` varchar(100) NOT NULL DEFAULT '' COMMENT '路径',
  `level` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '等级',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `role_name_index` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COMMENT='角色';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role`
--

LOCK TABLES `role` WRITE;
/*!40000 ALTER TABLE `role` DISABLE KEYS */;
INSERT INTO `role` VALUES (1,'sysadmin','系统管理员',0,'',1,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(2,'admin','管理员',1,'1-',2,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(3,'user','普通用户',2,'1-2-',3,'2022-04-01 00:00:00','2022-04-01 00:00:00');
/*!40000 ALTER TABLE `role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role_menu_list`
--

DROP TABLE IF EXISTS `role_menu_list`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `role_menu_list` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `role_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '角色id',
  `menu_list_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '菜单id',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COMMENT='角色和菜单关系';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role_menu_list`
--

LOCK TABLES `role_menu_list` WRITE;
/*!40000 ALTER TABLE `role_menu_list` DISABLE KEYS */;
INSERT INTO `role_menu_list` VALUES (1,1,1,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(2,1,2,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(3,1,3,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(4,1,4,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(5,1,5,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(6,1,6,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(7,3,1,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(8,3,3,'2022-04-01 00:00:00','2022-04-01 00:00:00'),(9,3,5,'2022-04-01 00:00:00','2022-04-01 00:00:00');
/*!40000 ALTER TABLE `role_menu_list` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `system`
--

DROP TABLE IF EXISTS `system`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `system` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `name` varchar(20) NOT NULL DEFAULT '' COMMENT '系统名称',
  `open_register` tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '开放注册 0/否 1/是',
  `register_quota` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '新默认流量 单位/byte',
  `register_expire_days` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '新用户默认过期天数 单位/天',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='系统设置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `system`
--

LOCK TABLES `system` WRITE;
/*!40000 ALTER TABLE `system` DISABLE KEYS */;
INSERT INTO `system` VALUES (2,'trojan-panel',1,0,0,'2022-04-01 00:00:00','2022-04-01 00:00:00');
/*!40000 ALTER TABLE `system` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` bigint(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `password` char(56) NOT NULL COMMENT '密码',
  `quota` bigint(20) NOT NULL DEFAULT '0' COMMENT '配额',
  `download` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '下载',
  `upload` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '上传',
  `username` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `pass` varchar(64) NOT NULL DEFAULT '' COMMENT '用户密码',
  `role_id` bigint(20) unsigned NOT NULL DEFAULT '3' COMMENT '角色id 1/系统管理员 3/普通用户',
  `deleted` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '是否禁用 0/正常 1/禁用',
  `expire_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '过期时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `password` (`password`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='账户';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'b4fc1369dd766eca295fb495b0938843becbac59fc5cb273b320aaa5',-1,0,0,'sysadmin','MTIzNDU2',1,0,32472115200000,'2022-04-01 00:00:00','2022-04-01 00:00:00');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-04-05  0:33:53
