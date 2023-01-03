-- MySQL dump 10.13  Distrib 5.7.29, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: trojan_panel_db
-- ------------------------------------------------------
-- Server version	5.7.29

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
-- Current Database: `trojan_panel_db`
--

CREATE DATABASE IF NOT EXISTS `trojan_panel_db` DEFAULT CHARACTER SET utf8mb4;

USE `trojan_panel_db`;

--
-- Table structure for table `account`
--

DROP TABLE IF EXISTS `account`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account` (
  `id` bigint(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `username` varchar(64) NOT NULL DEFAULT '' COMMENT '登录用户名',
  `pass` varchar(64) NOT NULL DEFAULT '' COMMENT '登录密码',
  `hash` varchar(64) NOT NULL DEFAULT '' COMMENT 'pass的hash',
  `role_id` bigint(20) unsigned NOT NULL DEFAULT '3' COMMENT '角色id 1/系统管理员 3/普通用户',
  `email` varchar(64) NOT NULL DEFAULT '' COMMENT '邮箱',
  `expire_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '过期时间',
  `deleted` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否禁用 0/正常 1/禁用',
  `quota` bigint(20) NOT NULL DEFAULT '0' COMMENT '配额 单位/byte',
  `download` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '下载 单位/byte',
  `upload` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '上传 单位/byte',
  `ip_limit` tinyint(2) unsigned NOT NULL DEFAULT '3' COMMENT '限制IP设备数',
  `upload_speed_limit` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '上传限速 单位/byte',
  `download_speed_limit` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '下载限速 单位/byte',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='账户';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `account`
--

LOCK TABLES `account` WRITE;
/*!40000 ALTER TABLE `account` DISABLE KEYS */;
INSERT INTO `account` VALUES (1,'sysadmin','tFjD2X1F6i9FfWp2GDU5Vbi1conuaChDKIYbw9zMFrqvMoSz','4366294571b8b267d9cf15b56660f0a70659568a86fc270a52fdc9e5',1,'',4078656000000,0,-1,0,0,3,0,0,'2022-04-01 00:00:00','2022-04-01 00:00:00');
/*!40000 ALTER TABLE `account` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `black_list`
--

DROP TABLE IF EXISTS `black_list`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `black_list` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `ip` varchar(64) NOT NULL DEFAULT '' COMMENT 'IP地址',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='黑名单';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `black_list`
--

LOCK TABLES `black_list` WRITE;
/*!40000 ALTER TABLE `black_list` DISABLE KEYS */;
/*!40000 ALTER TABLE `black_list` ENABLE KEYS */;
UNLOCK TABLES;

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
INSERT INTO `casbin_rule` VALUES ('p','sysadmin','/api/account/selectAccountById','GET','','',''),('p','sysadmin','/api/account/createAccount','POST','','',''),('p','sysadmin','/api/account/getAccountInfo','GET','','',''),('p','sysadmin','/api/account/selectAccountPage','GET','','',''),('p','sysadmin','/api/account/deleteAccountById','POST','','',''),('p','sysadmin','/api/account/updateAccountProfile','POST','','',''),('p','sysadmin','/api/account/updateAccountById','POST','','',''),('p','sysadmin','/api/account/logout','POST','','',''),('p','sysadmin','/api/account/clashSubscribe','GET','','',''),('p','sysadmin','/api/account/resetAccountDownloadAndUpload','POST','','',''),('p','sysadmin','/api/role/selectRoleList','GET','','',''),('p','sysadmin','/api/node/selectNodeById','GET','','',''),('p','sysadmin','/api/node/createNode','POST','','',''),('p','sysadmin','/api/node/selectNodePage','GET','','',''),('p','sysadmin','/api/node/deleteNodeById','POST','','',''),('p','sysadmin','/api/node/updateNodeById','POST','','',''),('p','sysadmin','/api/node/nodeQRCode','POST','','',''),('p','sysadmin','/api/node/nodeURL','POST','','',''),('p','sysadmin','/api/nodeType/selectNodeTypeList','GET','','',''),('p','sysadmin','/api/dashboard/panelGroup','GET','','',''),('p','sysadmin','/api/dashboard/trafficRank','GET','','',''),('p','sysadmin','/api/system/selectSystemByName','GET','','',''),('p','sysadmin','/api/system/updateSystemById','POST','','',''),('p','sysadmin','/api/system/uploadWebFile','POST','','',''),('p','sysadmin','/api/blackList/selectBlackListPage','GET','','',''),('p','sysadmin','/api/blackList/deleteBlackListByIp','POST','','',''),('p','sysadmin','/api/blackList/createBlackList','POST','','',''),('p','sysadmin','/api/emailRecord/selectEmailRecordPage','GET','','',''),('p','sysadmin','/api/nodeServer/selectNodeServerById','GET','','',''),('p','sysadmin','/api/nodeServer/createNodeServer','POST','','',''),('p','sysadmin','/api/nodeServer/selectNodeServerPage','GET','','',''),('p','sysadmin','/api/nodeServer/deleteNodeServerById','POST','','',''),('p','sysadmin','/api/nodeServer/updateNodeServerById','POST','','',''),('p','sysadmin','/api/nodeServer/selectNodeServerList','GET','','',''),('p','user','/api/account/getAccountInfo','GET','','',''),('p','user','/api/account/updateAccountProfile','POST','','',''),('p','user','/api/account/logout','POST','','',''),('p','user','/api/account/clashSubscribe','GET','','',''),('p','user','/api/node/selectNodePage','GET','','',''),('p','user','/api/node/nodeQRCode','POST','','',''),('p','user','/api/node/nodeURL','POST','','',''),('p','user','/api/nodeType/selectNodeTypeList','GET','','',''),('p','user','/api/dashboard/panelGroup','GET','','',''),('p','user','/api/dashboard/trafficRank','GET','','',''),('p','user','/api/nodeServer/selectNodeServerList','GET','','','');
/*!40000 ALTER TABLE `casbin_rule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `email_record`
--

DROP TABLE IF EXISTS `email_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `email_record` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `to_email` varchar(64) NOT NULL DEFAULT '' COMMENT '收件人邮箱',
  `subject` varchar(64) NOT NULL DEFAULT '' COMMENT '主题',
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '内容',
  `state` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '状态 0/未发送 1/发送成功 -1/发送失败',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='邮件发送记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `email_record`
--

LOCK TABLES `email_record` WRITE;
/*!40000 ALTER TABLE `email_record` DISABLE KEYS */;
/*!40000 ALTER TABLE `email_record` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `node`
--

DROP TABLE IF EXISTS `node`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `node` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `node_server_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '服务器id',
  `node_sub_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '节点分表id',
  `node_type_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '节点类型id',
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '名称',
  `node_server_ip` varchar(64) NOT NULL DEFAULT '' COMMENT 'IP地址',
  `domain` varchar(64) NOT NULL DEFAULT '' COMMENT '域名',
  `port` int(10) unsigned NOT NULL DEFAULT '443' COMMENT '端口',
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
-- Table structure for table `node_hysteria`
--

DROP TABLE IF EXISTS `node_hysteria`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `node_hysteria` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `protocol` varchar(32) NOT NULL DEFAULT 'udp' COMMENT '协议名称 udp/faketcp',
  `up_mbps` int(10) NOT NULL DEFAULT '100' COMMENT '单客户端最大上传速度 单位:Mbps',
  `down_mbps` int(10) NOT NULL DEFAULT '100' COMMENT '单客户端最大下载速度 单位:Mbps',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Hyeteria节点';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `node_hysteria`
--

LOCK TABLES `node_hysteria` WRITE;
/*!40000 ALTER TABLE `node_hysteria` DISABLE KEYS */;
/*!40000 ALTER TABLE `node_hysteria` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `node_server`
--

DROP TABLE IF EXISTS `node_server`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `node_server` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `ip` varchar(64) NOT NULL DEFAULT '' COMMENT '服务器IP',
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '服务器名称',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='服务器';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `node_server`
--

LOCK TABLES `node_server` WRITE;
/*!40000 ALTER TABLE `node_server` DISABLE KEYS */;
/*!40000 ALTER TABLE `node_server` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `node_trojan_go`
--

DROP TABLE IF EXISTS `node_trojan_go`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `node_trojan_go` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `sni` varchar(64) NOT NULL DEFAULT '' COMMENT 'sni',
  `mux_enable` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '是否开启多路复用 0/关闭 1/开启',
  `websocket_enable` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否开启websocket 0/否 1/是',
  `websocket_path` varchar(64) NOT NULL DEFAULT 'trojan-panel-websocket-path' COMMENT 'websocket路径',
  `websocket_host` varchar(64) NOT NULL DEFAULT '' COMMENT 'websocket host',
  `ss_enable` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否开启ss加密 0/否 1/是',
  `ss_method` varchar(32) NOT NULL DEFAULT 'AES-128-GCM' COMMENT 'ss加密方式',
  `ss_password` varchar(64) NOT NULL DEFAULT '' COMMENT 'ss密码',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='TrojanGO节点';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `node_trojan_go`
--

LOCK TABLES `node_trojan_go` WRITE;
/*!40000 ALTER TABLE `node_trojan_go` DISABLE KEYS */;
/*!40000 ALTER TABLE `node_trojan_go` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `node_type`
--

DROP TABLE IF EXISTS `node_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `node_type` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT '名称',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COMMENT='节点类型';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `node_type`
--

LOCK TABLES `node_type` WRITE;
/*!40000 ALTER TABLE `node_type` DISABLE KEYS */;
INSERT INTO `node_type` VALUES (1,'xray','2022-04-01 00:00:00','2022-04-01 00:00:00'),(2,'trojan-go','2022-04-01 00:00:00','2022-04-01 00:00:00'),(3,'hysteria','2022-04-01 00:00:00','2022-04-01 00:00:00'),(4,'naiveproxy','2022-04-01 00:00:00','2022-04-01 00:00:00');
/*!40000 ALTER TABLE `node_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `node_xray`
--

DROP TABLE IF EXISTS `node_xray`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `node_xray` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `protocol` varchar(32) NOT NULL DEFAULT '' COMMENT '协议名称',
  `settings` varchar(256) NOT NULL DEFAULT '' COMMENT 'settings',
  `stream_settings` varchar(256) NOT NULL DEFAULT '' COMMENT 'streamSettings',
  `tag` varchar(64) NOT NULL DEFAULT '' COMMENT 'tag',
  `sniffing` varchar(256) NOT NULL DEFAULT '' COMMENT 'sniffing',
  `allocate` varchar(256) NOT NULL DEFAULT '' COMMENT 'allocate',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Xray节点';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `node_xray`
--

LOCK TABLES `node_xray` WRITE;
/*!40000 ALTER TABLE `node_xray` DISABLE KEYS */;
/*!40000 ALTER TABLE `node_xray` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role`
--

DROP TABLE IF EXISTS `role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `role` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `name` varchar(16) NOT NULL DEFAULT '' COMMENT '名称',
  `desc` varchar(16) NOT NULL DEFAULT '' COMMENT '描述',
  `parent_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '父级id',
  `path` varchar(128) NOT NULL DEFAULT '' COMMENT '路径',
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
-- Table structure for table `system`
--

DROP TABLE IF EXISTS `system`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `system` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `name` varchar(16) NOT NULL DEFAULT '' COMMENT '系统名称',
  `account_config` varchar(512) NOT NULL DEFAULT '' COMMENT '用户设置',
  `email_config` varchar(512) NOT NULL DEFAULT '' COMMENT '系统邮箱设置',
  `template_config` varchar(5120) NOT NULL DEFAULT '' COMMENT '模板设置',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='系统设置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `system`
--

LOCK TABLES `system` WRITE;
/*!40000 ALTER TABLE `system` DISABLE KEYS */;
INSERT INTO `system` VALUES (1,'trojan-panel','{\"registerEnable\":1,\"registerQuota\":0,\"registerExpireDays\":0,\"resetDownloadAndUploadMonth\":0,\"trafficRankEnable\":1}','{\"expireWarnEnable\":0,\"expireWarnDay\":0,\"emailEnable\":0,\"emailHost\":\"\",\"emailPort\":0,\"emailUsername\":\"\",\"emailPassword\":\"\"}','{\"systemName\":\"Trojan Panel\",\"clashRule\":\"cnVsZXM6CiAgLSBSVUxFLVNFVCxhcHBsaWNhdGlvbnMsRElSRUNUCiAgLSBET01BSU4sY2xhc2gucmF6b3JkLnRvcCxESVJFQ1QKICAtIERPTUFJTix5YWNkLmhhaXNoYW4ubWUsRElSRUNUCiAgLSBSVUxFLVNFVCxwcml2YXRlLERJUkVDVAogIC0gUlVMRS1TRVQscmVqZWN0LFJFSkVDVAogIC0gUlVMRS1TRVQsaWNsb3VkLERJUkVDVAogIC0gUlVMRS1TRVQsYXBwbGUsRElSRUNUCiAgLSBSVUxFLVNFVCxnb29nbGUsRElSRUNUCiAgLSBSVUxFLVNFVCxwcm94eSxQUk9YWQogIC0gUlVMRS1TRVQsZGlyZWN0LERJUkVDVAogIC0gUlVMRS1TRVQsbGFuY2lkcixESVJFQ1QKICAtIFJVTEUtU0VULGNuY2lkcixESVJFQ1QKICAtIFJVTEUtU0VULHRlbGVncmFtY2lkcixQUk9YWQogIC0gR0VPSVAsLERJUkVDVAogIC0gR0VPSVAsQ04sRElSRUNUCiAgLSBNQVRDSCxQUk9YWQoKcnVsZS1wcm92aWRlcnM6CiAgcmVqZWN0OgogICAgdHlwZTogaHR0cAogICAgYmVoYXZpb3I6IGRvbWFpbgogICAgdXJsOiAiaHR0cHM6Ly9jZG4uanNkZWxpdnIubmV0L2doL0xveWFsc29sZGllci9jbGFzaC1ydWxlc0ByZWxlYXNlL3JlamVjdC50eHQiCiAgICBwYXRoOiAuL3J1bGVzZXQvcmVqZWN0LnlhbWwKICAgIGludGVydmFsOiA4NjQwMAoKICBpY2xvdWQ6CiAgICB0eXBlOiBodHRwCiAgICBiZWhhdmlvcjogZG9tYWluCiAgICB1cmw6ICJodHRwczovL2Nkbi5qc2RlbGl2ci5uZXQvZ2gvTG95YWxzb2xkaWVyL2NsYXNoLXJ1bGVzQHJlbGVhc2UvaWNsb3VkLnR4dCIKICAgIHBhdGg6IC4vcnVsZXNldC9pY2xvdWQueWFtbAogICAgaW50ZXJ2YWw6IDg2NDAwCgogIGFwcGxlOgogICAgdHlwZTogaHR0cAogICAgYmVoYXZpb3I6IGRvbWFpbgogICAgdXJsOiAiaHR0cHM6Ly9jZG4uanNkZWxpdnIubmV0L2doL0xveWFsc29sZGllci9jbGFzaC1ydWxlc0ByZWxlYXNlL2FwcGxlLnR4dCIKICAgIHBhdGg6IC4vcnVsZXNldC9hcHBsZS55YW1sCiAgICBpbnRlcnZhbDogODY0MDAKCiAgZ29vZ2xlOgogICAgdHlwZTogaHR0cAogICAgYmVoYXZpb3I6IGRvbWFpbgogICAgdXJsOiAiaHR0cHM6Ly9jZG4uanNkZWxpdnIubmV0L2doL0xveWFsc29sZGllci9jbGFzaC1ydWxlc0ByZWxlYXNlL2dvb2dsZS50eHQiCiAgICBwYXRoOiAuL3J1bGVzZXQvZ29vZ2xlLnlhbWwKICAgIGludGVydmFsOiA4NjQwMAoKICBwcm94eToKICAgIHR5cGU6IGh0dHAKICAgIGJlaGF2aW9yOiBkb21haW4KICAgIHVybDogImh0dHBzOi8vY2RuLmpzZGVsaXZyLm5ldC9naC9Mb3lhbHNvbGRpZXIvY2xhc2gtcnVsZXNAcmVsZWFzZS9wcm94eS50eHQiCiAgICBwYXRoOiAuL3J1bGVzZXQvcHJveHkueWFtbAogICAgaW50ZXJ2YWw6IDg2NDAwCgogIGRpcmVjdDoKICAgIHR5cGU6IGh0dHAKICAgIGJlaGF2aW9yOiBkb21haW4KICAgIHVybDogImh0dHBzOi8vY2RuLmpzZGVsaXZyLm5ldC9naC9Mb3lhbHNvbGRpZXIvY2xhc2gtcnVsZXNAcmVsZWFzZS9kaXJlY3QudHh0IgogICAgcGF0aDogLi9ydWxlc2V0L2RpcmVjdC55YW1sCiAgICBpbnRlcnZhbDogODY0MDAKCiAgcHJpdmF0ZToKICAgIHR5cGU6IGh0dHAKICAgIGJlaGF2aW9yOiBkb21haW4KICAgIHVybDogImh0dHBzOi8vY2RuLmpzZGVsaXZyLm5ldC9naC9Mb3lhbHNvbGRpZXIvY2xhc2gtcnVsZXNAcmVsZWFzZS9wcml2YXRlLnR4dCIKICAgIHBhdGg6IC4vcnVsZXNldC9wcml2YXRlLnlhbWwKICAgIGludGVydmFsOiA4NjQwMAoKICBnZnc6CiAgICB0eXBlOiBodHRwCiAgICBiZWhhdmlvcjogZG9tYWluCiAgICB1cmw6ICJodHRwczovL2Nkbi5qc2RlbGl2ci5uZXQvZ2gvTG95YWxzb2xkaWVyL2NsYXNoLXJ1bGVzQHJlbGVhc2UvZ2Z3LnR4dCIKICAgIHBhdGg6IC4vcnVsZXNldC9nZncueWFtbAogICAgaW50ZXJ2YWw6IDg2NDAwCgogIGdyZWF0ZmlyZToKICAgIHR5cGU6IGh0dHAKICAgIGJlaGF2aW9yOiBkb21haW4KICAgIHVybDogImh0dHBzOi8vY2RuLmpzZGVsaXZyLm5ldC9naC9Mb3lhbHNvbGRpZXIvY2xhc2gtcnVsZXNAcmVsZWFzZS9ncmVhdGZpcmUudHh0IgogICAgcGF0aDogLi9ydWxlc2V0L2dyZWF0ZmlyZS55YW1sCiAgICBpbnRlcnZhbDogODY0MDAKCiAgdGxkLW5vdC1jbjoKICAgIHR5cGU6IGh0dHAKICAgIGJlaGF2aW9yOiBkb21haW4KICAgIHVybDogImh0dHBzOi8vY2RuLmpzZGVsaXZyLm5ldC9naC9Mb3lhbHNvbGRpZXIvY2xhc2gtcnVsZXNAcmVsZWFzZS90bGQtbm90LWNuLnR4dCIKICAgIHBhdGg6IC4vcnVsZXNldC90bGQtbm90LWNuLnlhbWwKICAgIGludGVydmFsOiA4NjQwMAoKICB0ZWxlZ3JhbWNpZHI6CiAgICB0eXBlOiBodHRwCiAgICBiZWhhdmlvcjogaXBjaWRyCiAgICB1cmw6ICJodHRwczovL2Nkbi5qc2RlbGl2ci5uZXQvZ2gvTG95YWxzb2xkaWVyL2NsYXNoLXJ1bGVzQHJlbGVhc2UvdGVsZWdyYW1jaWRyLnR4dCIKICAgIHBhdGg6IC4vcnVsZXNldC90ZWxlZ3JhbWNpZHIueWFtbAogICAgaW50ZXJ2YWw6IDg2NDAwCgogIGNuY2lkcjoKICAgIHR5cGU6IGh0dHAKICAgIGJlaGF2aW9yOiBpcGNpZHIKICAgIHVybDogImh0dHBzOi8vY2RuLmpzZGVsaXZyLm5ldC9naC9Mb3lhbHNvbGRpZXIvY2xhc2gtcnVsZXNAcmVsZWFzZS9jbmNpZHIudHh0IgogICAgcGF0aDogLi9ydWxlc2V0L2NuY2lkci55YW1sCiAgICBpbnRlcnZhbDogODY0MDAKCiAgbGFuY2lkcjoKICAgIHR5cGU6IGh0dHAKICAgIGJlaGF2aW9yOiBpcGNpZHIKICAgIHVybDogImh0dHBzOi8vY2RuLmpzZGVsaXZyLm5ldC9naC9Mb3lhbHNvbGRpZXIvY2xhc2gtcnVsZXNAcmVsZWFzZS9sYW5jaWRyLnR4dCIKICAgIHBhdGg6IC4vcnVsZXNldC9sYW5jaWRyLnlhbWwKICAgIGludGVydmFsOiA4NjQwMAoKICBhcHBsaWNhdGlvbnM6CiAgICB0eXBlOiBodHRwCiAgICBiZWhhdmlvcjogY2xhc3NpY2FsCiAgICB1cmw6ICJodHRwczovL2Nkbi5qc2RlbGl2ci5uZXQvZ2gvTG95YWxzb2xkaWVyL2NsYXNoLXJ1bGVzQHJlbGVhc2UvYXBwbGljYXRpb25zLnR4dCIKICAgIHBhdGg6IC4vcnVsZXNldC9hcHBsaWNhdGlvbnMueWFtbAogICAgaW50ZXJ2YWw6IDg2NDAw\"}','2022-04-01 00:00:00','2022-04-01 00:00:00');
/*!40000 ALTER TABLE `system` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-01-03 15:09:10
