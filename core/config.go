package core

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"os"
	"runtime"
	"strconv"
	"trojan-panel/module/constant"
	"trojan-panel/util"
)

var (
	host           string
	user           string
	password       string
	port           string
	redisHost      string
	redisPort      string
	redisPassword  string
	redisDb        string
	redisMaxIdle   string
	redisMaxActive string
	redisWait      string
	serverPort     string
	version        bool
)

func init() {
	flag.StringVar(&host, "host", "localhost", "数据库地址")
	flag.StringVar(&user, "user", "root", "数据库用户名")
	flag.StringVar(&password, "password", "123456", "数据库密码")
	flag.StringVar(&port, "port", "3306", "数据库端口")
	flag.StringVar(&redisHost, "redisHost", "127.0.0.1", "Redis地址")
	flag.StringVar(&redisPort, "redisPort", "6379", "Redis端口")
	flag.StringVar(&redisPassword, "redisPassword", "123456", "Redis密码")
	flag.StringVar(&redisDb, "redisDb", "0", "Redis默认数据库")
	flag.StringVar(&redisMaxIdle, "redisMaxIdle", strconv.FormatInt(int64(runtime.NumCPU()*2), 10), "Redis最大空闲连接数")
	flag.StringVar(&redisMaxActive, "redisMaxActive", strconv.FormatInt(int64(runtime.NumCPU()*2+2), 10), "Redis最大连接数")
	flag.StringVar(&redisWait, "redisWait", "true", "Redis是否等待")
	flag.StringVar(&serverPort, "serverPort", "8081", "服务端口")
	flag.BoolVar(&version, "version", false, "打印版本信息")
	flag.Usage = usage
	flag.Parse()
	if version {
		_, _ = fmt.Fprint(os.Stdout, constant.TrojanPanelVersion)
		os.Exit(0)
	}

	logPath := constant.LogPath
	if !util.Exists(logPath) {
		if err := os.Mkdir(logPath, os.ModePerm); err != nil {
			logrus.Errorf("创建logs文件夹异常 err: %v", err)
			panic(err)
		}
	}

	webFilePath := constant.WebFilePath
	if !util.Exists(webFilePath) {
		if err := os.Mkdir(webFilePath, os.ModePerm); err != nil {
			logrus.Errorf("创建webfile文件夹异常 err: %v", err)
			panic(err)
		}
	}

	configPath := constant.ConfigPath
	if !util.Exists(configPath) {
		if err := os.Mkdir(configPath, os.ModePerm); err != nil {
			logrus.Errorf("创建config文件夹异常 err: %v", err)
			panic(err)
		}
	}

	templatePath := constant.TemplatePath
	if !util.Exists(templatePath) {
		if err := os.Mkdir(templatePath, os.ModePerm); err != nil {
			logrus.Errorf("创建config/template文件夹异常 err: %v", err)
			panic(err)
		}
	}

	// 创建Logo
	logoImagePath := constant.LogoImagePath
	if !util.Exists(logoImagePath) {
		if err := util.DownloadFile(constant.LogoImageUrl, logoImagePath); err != nil {
			logrus.Errorf("创建logo.png文件异常 err: %v", err)
		}
	}

	ExportPath := constant.ExportPath
	if !util.Exists(ExportPath) {
		if err := os.Mkdir(ExportPath, os.ModePerm); err != nil {
			logrus.Errorf("创建config/export文件夹异常 err: %v", err)
			panic(err)
		}
	}

	// 创建AccountTemplate模板
	exportAccountTemplate := constant.ExportAccountTemplate
	if !util.Exists(exportAccountTemplate) {
		var accountTemplate []map[string]any
		accountTemplate = append(accountTemplate, map[string]any{"username": "example", "pass": "example", "hash": "example", "role_id": 3, "email": "test@example.com", "expire_time": int64(4078656000000), "deleted": 0, "quota": -1, "download": 0, "upload": 0})
		if err := util.ExportJson(exportAccountTemplate, accountTemplate); err != nil {
			logrus.Errorf("创建AccountTemplate.json异常 err: %v", err)
			panic(err)
		}
	}

	// 创建NodeServerTemplate模板
	exportNodeServerTemplate := constant.ExportNodeServerTemplate
	if !util.Exists(exportNodeServerTemplate) {
		var nodeServerTemplate []map[string]any
		nodeServerTemplate = append(nodeServerTemplate, map[string]any{"ip": "127.0.0.1", "name": "example", "grpc_port": 8100})
		if err := util.ExportJson(exportNodeServerTemplate, nodeServerTemplate); err != nil {
			logrus.Errorf("创建NodeServerTemplate.json异常 err: %v", err)
			panic(err)
		}
	}

	configFilePath := constant.ConfigFilePath
	if !util.Exists(configFilePath) {
		file, err := os.Create(configFilePath)
		if err != nil {
			logrus.Errorf("创建config.ini文件异常 err: %v", err)
			panic(err)
		}
		defer file.Close()

		_, err = file.WriteString(fmt.Sprintf(
			`[mysql]
host=%s
user=%s
password=%s
port=%s
[log]
filename=logs/trojan-panel.log
max_size=1
max_backups=5
max_age=30
compress=true
[redis]
host=%s
port=%s
password=%s
db=%s
max_idle=%s
max_active=%s
wait=%s
[server]
port=%s
`, host, user, password, port, redisHost, redisPort, redisPassword, redisDb,
			redisMaxIdle, redisMaxActive, redisWait, serverPort))
		if err != nil {
			logrus.Errorf("config.ini文件写入异常 err: %v", err)
			panic(err)
		}

	}

	rbacModelConfigPath := constant.RbacModelFilePath
	if !util.Exists(rbacModelConfigPath) {
		file, err := os.Create(rbacModelConfigPath)
		if err != nil {
			logrus.Errorf("创建rbac_model.conf文件异常 err: %v", err)
			panic(err)
		}
		defer file.Close()

		_, err = file.WriteString(
			`[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`)
		if err != nil {
			logrus.Errorf("rbac_model.conf文件写入异常 err: %v", err)
			panic(err)
		}
	}

	clashRuleFilePath := constant.ClashRuleFilePath
	if !util.Exists(clashRuleFilePath) {
		file, err := os.Create(clashRuleFilePath)
		if err != nil {
			logrus.Errorf("创建clash-rule.yaml文件异常 err: %v", err)
			panic(err)
		}
		defer file.Close()

		_, err = file.WriteString(constant.ClashRules)
		if err != nil {
			logrus.Errorf("clash-rule.yaml文件写入异常 err: %v", err)
			panic(err)
		}
	}

	xrayTemplateFilePath := constant.XrayTemplateFilePath
	if !util.Exists(xrayTemplateFilePath) {
		file, err := os.Create(xrayTemplateFilePath)
		if err != nil {
			logrus.Errorf("创建template-xray.json文件异常 err: %v", err)
			panic(err)
		}
		defer file.Close()

		_, err = file.WriteString(`{
  "log": {
    "loglevel": "warning"
  },
  "inbounds": [],
  "outbounds": [
    {
      "protocol": "freedom"
    }
  ],
  "api": {
    "tag": "api",
    "services": [
      "HandlerService",
      "LoggerService",
      "StatsService"
    ]
  },
  "routing": {
    "rules": [
      {
        "inboundTag": [
          "api"
        ],
        "outboundTag": "api",
        "type": "field"
      }
    ]
  },
  "stats": {},
  "policy": {
    "levels": {
      "0": {
        "statsUserUplink": true,
        "statsUserDownlink": true
      }
    },
    "system": {
      "statsInboundUplink": true,
      "statsInboundDownlink": true
    }
  }
}`)
		if err != nil {
			logrus.Errorf("template-xray.json文件写入异常 err: %v", err)
			panic(err)
		}
	}
}

func usage() {
	_, _ = fmt.Fprintln(os.Stdout, `trojan panel help
Usage: trojan-panel [-host] [-password] [-port] [-redisHost] [-redisPort] [-redisPassword] [-redisDb] [-redisMaxIdle] [-redisMaxActive] [-redisWait] [-h] [-version]`)
	flag.PrintDefaults()
}

var Config = new(AppConfig)

// InitConfig 初始化全局配置文件
func InitConfig() {
	if err := ini.MapTo(Config, constant.ConfigFilePath); err != nil {
		logrus.Errorf("配置文件加载失败 err: %v", err)
		panic(err)
	}
}

type AppConfig struct {
	MySQLConfig  `ini:"mysql"`
	LogConfig    `ini:"log"`
	RedisConfig  `ini:"redis"`
	ServerConfig `ini:"server"`
}

type MySQLConfig struct {
	Host     string `ini:"host"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	Port     int    `ini:"port"`
}

type LogConfig struct {
	FileName   string `ini:"filename"`    // 日志文件位置
	MaxSize    int    `ini:"max_size"`    // 单文件最大容量,单位是MB
	MaxBackups int    `ini:"max_backups"` // 最大保留过期文件个数
	MaxAge     int    `ini:"max_age"`     // 保留过期文件的最大时间间隔,单位是天
	Compress   bool   `ini:"compress"`    // 是否需要压缩滚动日志, 使用的 gzip 压缩
}

type RedisConfig struct {
	Host      string `ini:"host"`
	Port      int    `ini:"port"`
	Password  string `ini:"password"`
	Db        int    `ini:"db"`
	MaxIdle   int    `ini:"max_idle"`
	MaxActive int    `ini:"max_active"`
	Wait      bool   `ini:"wait"`
}

type ServerConfig struct {
	Port int `ini:"port"`
}
