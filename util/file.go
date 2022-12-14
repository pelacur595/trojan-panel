package util

import (
	"archive/zip"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"trojan-panel/module/constant"
)

func DownloadFile(url string, fileName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(fileName, data, 0644); err != nil {
		return err
	}
	return nil
}

// 解压
func Unzip(src string, dest string) error {
	// 打开读取压缩文件
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// 遍历压缩文件内的文件，写入磁盘
	for _, f := range r.File {
		filePath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: 非法的文件路径", filePath)
		}

		// 如果是目录，就创建目录
		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}

		rc.Close()
		outFile.Close()
	}
	return nil
}

// 删除文件夹内的子文件包括目录
func RemoveSubFile(filePath string) error {
	dir, err := ioutil.ReadDir(filePath)
	if err != nil {
		return err
	}
	for _, d := range dir {
		if err := os.RemoveAll(path.Join([]string{filePath, d.Name()}...)); err != nil {
			return err
		}
	}
	return nil
}

// 判断文件或者文件夹是否存在
func Exists(path string) bool {
	// 获取文件信息
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 初始化文件/文件夹
func InitFile() {
	logPath := constant.LogPath
	if !Exists(logPath) {
		if err := os.Mkdir(logPath, os.ModePerm); err != nil {
			logrus.Errorf("创建logs文件夹异常 err: %v", err)
			panic(err)
		}
	}
	webFilePath := constant.WebFilePath
	if !Exists(webFilePath) {
		if err := os.Mkdir(webFilePath, os.ModePerm); err != nil {
			logrus.Errorf("创建webfile文件夹异常 err: %v", err)
			panic(err)
		}
	}
	configPath := constant.ConfigPath
	if !Exists(configPath) {
		if err := os.Mkdir(configPath, os.ModePerm); err != nil {
			logrus.Errorf("创建config文件夹异常 err: %v", err)
			panic(err)
		}
	}

	configFilePath := constant.ConfigFilePath
	if !Exists(configFilePath) {
		file, err := os.Create(configFilePath)
		if err != nil {
			logrus.Errorf("创建config.ini文件异常 err: %v", err)
			panic(err)
		}
		defer file.Close()

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
			version        bool
		)
		flag.StringVar(&host, "host", "localhost", "数据库地址")
		flag.StringVar(&user, "user", "root", "数据库用户名")
		flag.StringVar(&password, "password", "123456", "数据库密码")
		flag.StringVar(&port, "port", "3306", "数据库端口")
		flag.StringVar(&redisHost, "redisHost", "127.0.0.1", "Redis地址")
		flag.StringVar(&redisPort, "redisPort", "6379", "Redis端口")
		flag.StringVar(&redisPassword, "redisPassword", "123456", "Redis密码")
		flag.StringVar(&redisDb, "redisDb", "0", "Redis默认数据库")
		flag.StringVar(&redisMaxIdle, "redisMaxIdle", "2", "Redis最大空闲连接数")
		flag.StringVar(&redisMaxActive, "redisMaxActive", "2", "Redis最大连接数")
		flag.StringVar(&redisWait, "redisWait", "true", "Redis是否等待")
		flag.BoolVar(&version, "version", false, "print trojan panel version")
		flag.Parse()
		if version {
			println(constant.TrojanPanelVersion)
			os.Exit(0)
		}
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
`, host, user, password, port, redisHost, redisPort, redisPassword, redisDb,
			redisMaxIdle, redisMaxIdle, redisWait))
		if err != nil {
			logrus.Errorf("config.ini文件写入异常 err: %v", err)
			panic(err)
		}
		flag.Usage = usage
	}

	rbacModelConfigPath := constant.RbacModelFilePath
	if !Exists(rbacModelConfigPath) {
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
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `trojan panel help
Usage: trojanpanel [-host] [-password] [-port] [-redisHost] [-redisPort] [-redisPassword] [-redisDb] [-redisMaxIdle] [-redisMaxActive] [-redisWait] [-h] [-version]

Options:
	-host			database host
	-user			database user
	-password		database password
	-port			database port
	-redisHost		redis redisHost
	-redisPort		redis redisPort
	-redisPassword	redis redisPassword
	-redisDb		redis redisDb
	-redisMaxIdle   redis redisMaxIdle
	-redisMaxActive	redis redisMaxActive
	-redisWait		redis redisWait
	-h				help
	-version        trojan panel version
`)
}
