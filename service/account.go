package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"mime/multipart"
	"sync"
	"time"
	"trojan-panel/core"
	"trojan-panel/dao"
	"trojan-panel/dao/redis"
	"trojan-panel/module"
	"trojan-panel/module/bo"
	"trojan-panel/module/constant"
	"trojan-panel/module/dto"
	"trojan-panel/module/vo"
	"trojan-panel/util"
)

func CreateAccount(accountCreateDto dto.AccountCreateDto) error {
	count, err := dao.CountAccountByUsername(accountCreateDto.Username)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(constant.UsernameExist)
	}
	toByte := util.ToByte(*accountCreateDto.Quota)
	account := module.Account{
		Username:   accountCreateDto.Username,
		Pass:       accountCreateDto.Pass,
		RoleId:     accountCreateDto.RoleId,
		Email:      accountCreateDto.Email,
		ExpireTime: accountCreateDto.ExpireTime,
		Deleted:    accountCreateDto.Deleted,
		Quota:      &toByte,
		//IpLimit:            accountCreateDto.IpLimit,
		//DownloadSpeedLimit: accountCreateDto.DownloadSpeedLimit,
		//UploadSpeedLimit:   accountCreateDto.UploadSpeedLimit,
	}
	if err = dao.CreateAccount(&account); err != nil {
		return err
	}
	if account.Deleted != nil && *account.Deleted == 1 {
		if err = PullAccountWhiteOrBlackByUsername([]string{*account.Username}, true); err != nil {
			return err
		}
	} else if *account.ExpireTime <= util.NowMilli() {
		if err = DisableAccount([]string{*account.Username}); err != nil {
			return err
		}
	}
	return nil
}
func SelectAccountById(id *uint) (*module.Account, error) {
	return dao.SelectAccountById(id)
}
func CountAccountByUsername(username *string) (int, error) {
	return dao.CountAccountByUsername(username)
}
func SelectAccountPage(username *string, deleted *uint, pageNum *uint, pageSize *uint) (*vo.AccountPageVo, error) {
	return dao.SelectAccountPage(username, deleted, pageNum, pageSize)
}
func DeleteAccountById(token string, id *uint) error {
	var mutex sync.Mutex
	defer mutex.Unlock()
	if mutex.TryLock() {
		password, err := dao.SelectConnectPassword(id, nil)
		if err != nil {
			return err
		}
		if err = RemoveAccount(token, password); err != nil {
			return err
		}
		if err = dao.DeleteAccountById(id); err != nil {
			return err
		}
	}
	return nil
}

func SelectAccountByUsername(username *string) (*module.Account, error) {
	return dao.SelectAccountByUsername(username)
}

func UpdateAccountPass(token string, oldPass *string, newPass *string, username *string) error {
	var mutex sync.Mutex
	defer mutex.Unlock()
	if mutex.TryLock() {
		account, err := SelectAccountByUsername(username)
		if err != nil || !util.Sha1Match(*account.Pass, fmt.Sprintf("%s%s", *username, *oldPass)) {
			return errors.New(constant.OriPassError)
		}

		if err = RemoveAccount(token, *account.Pass); err != nil {
			return err
		}

		if err := dao.UpdateAccountPass(oldPass, newPass, username); err != nil {
			return err
		}
	}
	return nil
}

func UpdateAccountProperty(token string, oldUsername *string, pass *string, username *string, email *string) error {
	var mutex sync.Mutex
	defer mutex.Unlock()
	if mutex.TryLock() {
		account, err := SelectAccountByUsername(oldUsername)
		if err != nil || !util.Sha1Match(*account.Pass, fmt.Sprintf("%s%s", *oldUsername, *pass)) {
			return errors.New(constant.OriPassError)
		}

		if pass != nil && *pass != "" && username != nil && *username != "" {
			if err = RemoveAccount(token, *account.Pass); err != nil {
				return err
			}
		}

		if err := dao.UpdateAccountProperty(oldUsername, pass, username, email); err != nil {
			return err
		}
	}
	return nil
}

// GetAccountInfo 获取当前请求账户信息
func GetAccountInfo(c *gin.Context) (*vo.AccountInfo, error) {
	accountVo := util.GetCurrentAccount(c)
	roles, err := dao.SelectRoleNameByParentId(&accountVo.RoleId, true)
	if err != nil {
		return nil, err
	}
	userInfo := vo.AccountInfo{
		Id:       accountVo.Id,
		Username: accountVo.Username,
		Roles:    roles,
	}
	return &userInfo, nil
}

func UpdateAccountById(token string, account *module.Account) error {
	var mutex sync.Mutex
	defer mutex.Unlock()
	if mutex.TryLock() {
		if account.Pass != nil && *account.Pass != "" {
			password, err := dao.SelectConnectPassword(account.Id, nil)
			if err != nil {
				return err
			}
			if err = RemoveAccount(token, password); err != nil {
				return err
			}
		}
		if err := dao.UpdateAccountById(account); err != nil {
			return err
		}
		if account.Deleted != nil && *account.Deleted == 1 {
			if err := PullAccountWhiteOrBlackByUsername([]string{*account.Username}, true); err != nil {
				return err
			}
		} else if *account.ExpireTime <= util.NowMilli() {
			if err := DisableAccount([]string{*account.Username}); err != nil {
				return err
			}
		}
	}
	return nil
}

func Register(accountRegisterDto dto.AccountRegisterDto) error {
	name := constant.SystemName
	systemVo, err := SelectSystemByName(&name)
	if err != nil {
		return err
	}
	if systemVo.RegisterEnable == 0 {
		return errors.New(constant.AccountRegisterClosed)
	}

	count, err := dao.CountAccountByUsername(accountRegisterDto.Username)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(constant.UsernameExist)
	}
	u := constant.USER
	milli := util.DayToMilli(systemVo.RegisterExpireDays)
	registerQuota := util.ToByte(systemVo.RegisterQuota)
	account := module.Account{
		Quota:      &registerQuota,
		Username:   accountRegisterDto.Username,
		Pass:       accountRegisterDto.Pass,
		RoleId:     &u,
		Deleted:    new(uint),
		ExpireTime: &milli,
	}
	if err = dao.CreateAccount(&account); err != nil {
		return err
	}
	return nil
}

// PullAccountWhiteOrBlackByUsername 拉白或者拉黑用户 此操作会清空用户流量
func PullAccountWhiteOrBlackByUsername(usernames []string, isBlack bool) error {
	if len(usernames) > 0 {
		var deleted uint
		if isBlack {
			deleted = 1
		} else {
			deleted = 0
		}
		if err := dao.UpdateAccountQuotaOrDownloadOrUploadOrDeletedByUsernames(usernames, new(int), new(uint), new(uint), &deleted); err != nil {
			return err
		}
	}
	return nil
}

// DisableAccount 清空流量/禁用用户连接节点
func DisableAccount(usernames []string) error {
	if len(usernames) > 0 {
		if err := dao.UpdateAccountQuotaOrDownloadOrUploadOrDeletedByUsernames(usernames, new(int), new(uint), new(uint), nil); err != nil {
			return err
		}
	}
	return nil
}

// CronScanAccounts 定时任务：扫描无效用户
func CronScanAccounts() {
	usernames, err := dao.SelectAccountUsernameByDeletedOrExpireTime()
	if err != nil {
		return
	}

	if len(usernames) > 0 {
		if err = DisableAccount(usernames); err != nil {
			logrus.Errorf("定时扫描用户任务禁用用户异常 usernames: %s error: %v", usernames, err)
		}
		logrus.Infof("定时扫描用户任务禁用用户 usernames: %s", usernames)
	}
}

// CronScanAccountExpireWarn 定时任务：到期警告
func CronScanAccountExpireWarn() {
	systemName := constant.SystemName
	systemVo, err := SelectSystemByName(&systemName)
	if err != nil {
		return
	}
	if systemVo.EmailEnable == 0 || systemVo.ExpireWarnEnable == 0 {
		return
	}
	expireWarnDay := systemVo.ExpireWarnDay
	accounts, err := dao.SelectAccountsByExpireTime(util.DayToMilli(expireWarnDay))
	if err != nil {
		logrus.Errorln(err.Error())
		return
	}
	if len(accounts) > 0 {
		for _, account := range accounts {
			if account.Email != nil && *account.Email != "" {
				// 发送到期邮件
				emailDto := dto.SendEmailDto{
					FromEmailName: systemVo.SystemName,
					ToEmails:      []string{*account.Email},
					Subject:       "账号到期提醒",
					Content:       fmt.Sprintf("您的账户: %s,还有%d天到期,请及时续期", *account.Username, expireWarnDay),
				}
				if err = SendEmail(&emailDto); err != nil {
					logrus.Errorln(fmt.Sprintf("到期警告邮件发送失败 err: %v", err))
				}
			}
		}
	}
}

// CronResetDownloadAndUploadMonth 定时任务：每月重设除管理员之外的用户下载和上传流量
func CronResetDownloadAndUploadMonth() {
	name := constant.SystemName
	systemConfig, err := SelectSystemByName(&name)
	if err != nil {
		logrus.Errorf("每月重设除管理员之外的用户下载和上传流量 查询系统设置异常 error: %v", err)
		return
	}
	if systemConfig.ResetDownloadAndUploadMonth == 1 {
		roleIds := []uint{constant.USER}
		if err := dao.ResetAccountDownloadAndUpload(nil, &roleIds); err != nil {
			logrus.Errorf("每月重设除管理员之外的用户下载和上传流量异常 roleIds: %v error: %v", roleIds, err)
		}
	}
}

func RemoveAccount(token string, password string) error {
	nodes, err := dao.SelectNodesIpGrpcPortDistinct()
	if err != nil {
		return err
	}
	for _, node := range nodes {
		removeDto := core.AccountRemoveDto{
			Password: password,
		}
		_ = core.RemoveAccount(token, *node.NodeServerIp, *node.NodeServerGrpcPort, &removeDto)
	}
	return nil
}

func SelectConnectPassword(id *uint, username *string) (string, error) {
	return dao.SelectConnectPassword(id, username)
}

// ResetAccountDownloadAndUpload 重设下载和上传流量
func ResetAccountDownloadAndUpload(id *uint, roleIds *[]uint) error {
	return dao.ResetAccountDownloadAndUpload(id, roleIds)
}

// SubscribeClash
/**
Clash for windows 参考文档：
1. https://docs.cfw.lbyczf.com/contents/urlscheme.html
2. https://github.com/crossutility/Quantumult/blob/master/extra-subscription-feature.md
3. https://github.com/Dreamacro/clash/wiki/Configuration
*/
func SubscribeClash(pass string) (*module.Account, string, []byte, vo.SystemVo, error) {
	account, err := dao.SelectAccountClashSubscribe(pass)
	if err != nil {
		return nil, "", []byte{}, vo.SystemVo{}, err
	}
	nodes, err := dao.SelectNodesIpAndPort()
	if err != nil {
		return nil, "", []byte{}, vo.SystemVo{}, err
	}
	var nodeOneVos []vo.NodeOneVo
	for _, item := range nodes {
		nodeOneVo, err := SelectNodeById(item.Id)
		if err != nil {
			continue
		}
		nodeOneVos = append(nodeOneVos, *nodeOneVo)
	}

	userInfo := fmt.Sprintf("upload=%d; download=%d; total=%d; expire=%d",
		*account.Upload,
		*account.Download,
		*account.Quota,
		*account.ExpireTime/1000)

	clashConfig := bo.ClashConfig{}
	var ClashConfigInterface []interface{}
	var proxies []string
	for _, item := range nodeOneVos {
		if item.NodeTypeId == constant.Xray {
			nodeXray, err := SelectNodeXrayById(&item.NodeSubId)
			if err != nil {
				return nil, "", []byte{}, vo.SystemVo{}, err
			}

			streamSettings := bo.StreamSettings{}
			if nodeXray.StreamSettings != nil && *nodeXray.StreamSettings != "" {
				if err = json.Unmarshal([]byte(*nodeXray.StreamSettings), &streamSettings); err != nil {
					logrus.Errorln(fmt.Sprintf("SystemVo JSON反转失败 err: %v", err))
					return nil, "", []byte{}, vo.SystemVo{}, errors.New(constant.SysError)
				}
			}
			settings := bo.Settings{}
			if nodeXray.Settings != nil && *nodeXray.Settings != "" {
				if err = json.Unmarshal([]byte(*nodeXray.Settings), &settings); err != nil {
					logrus.Errorln(fmt.Sprintf("SystemVo JSON反转失败 err: %v", err))
					return nil, "", []byte{}, vo.SystemVo{}, errors.New(constant.SysError)
				}
			}
			switch *nodeXray.Protocol {
			case constant.ProtocolVless:
				vless := bo.Vless{}
				vless.Name = item.Name
				vless.Type = constant.ClashVless
				vless.Server = item.Domain
				vless.Port = item.Port
				vless.Uuid = util.GenerateUUID(pass)
				vless.Network = streamSettings.Network
				vless.Tls = true
				vless.Udp = true
				vless.Flow = item.XrayFlow
				if streamSettings.Security == "tls" {
					vless.ClientFingerprint = streamSettings.TlsSettings.Fingerprint
					vless.SkipCertVerify = streamSettings.TlsSettings.AllowInsecure
					vless.ServerName = streamSettings.TlsSettings.ServerName
				} else if streamSettings.Security == "reality" {
					if len(streamSettings.RealitySettings.ServerNames) > 0 {
						vless.ServerName = streamSettings.RealitySettings.ServerNames[0]
					}
					if len(streamSettings.RealitySettings.ShortIds) > 0 {
						vless.RealityOpts.ShortId = streamSettings.RealitySettings.ShortIds[0]
					}
					vless.RealityOpts.PublicKey = item.RealityPbk
					vless.ClientFingerprint = streamSettings.RealitySettings.Fingerprint
				}
				if streamSettings.Network == "ws" {
					vless.WsOpts.Path = streamSettings.WsSettings.Path
					vless.WsOpts.Headers.Host = streamSettings.WsSettings.Headers.Host
				}
				ClashConfigInterface = append(ClashConfigInterface, vless)
				proxies = append(proxies, item.Name)
			case constant.ProtocolVmess:
				vmess := bo.Vmess{}
				vmess.Name = item.Name
				vmess.Type = constant.ClashVmess
				vmess.Server = item.Domain
				vmess.Port = item.Port
				vmess.Uuid = util.GenerateUUID(pass)
				vmess.AlterId = 0
				if settings.Encryption != "none" {
					vmess.Cipher = "auto"
				} else {
					vmess.Cipher = "none"
				}
				vmess.Udp = true
				vmess.Tls = true
				vmess.Network = streamSettings.Network
				if streamSettings.Security == "tls" {
					vmess.ClientFingerprint = streamSettings.TlsSettings.Fingerprint
					vmess.SkipCertVerify = streamSettings.TlsSettings.AllowInsecure
					vmess.ServerName = streamSettings.TlsSettings.ServerName
				}
				if streamSettings.Network == "ws" {
					vmess.WsOpts.Path = streamSettings.WsSettings.Path
					vmess.WsOpts.Headers.Host = streamSettings.WsSettings.Headers.Host
				}
				ClashConfigInterface = append(ClashConfigInterface, vmess)
				proxies = append(proxies, item.Name)
			case constant.ProtocolTrojan:
				trojan := bo.Trojan{}
				trojan.Name = item.Name
				trojan.Type = constant.ClashTrojan
				trojan.Server = item.Domain
				trojan.Port = item.Port
				trojan.Password = pass
				trojan.Udp = true
				if streamSettings.Security == "tls" {
					trojan.ClientFingerprint = streamSettings.TlsSettings.Fingerprint
					trojan.Sni = streamSettings.TlsSettings.ServerName
					trojan.Alpn = streamSettings.TlsSettings.Alpn
					trojan.SkipCertVerify = streamSettings.TlsSettings.AllowInsecure
				}
				if streamSettings.Network == "ws" {
					trojan.WsOpts.Path = streamSettings.WsSettings.Path
					trojan.WsOpts.Headers.Host = streamSettings.WsSettings.Headers.Host
				}
				ClashConfigInterface = append(ClashConfigInterface, trojan)
				proxies = append(proxies, item.Name)
			case constant.ProtocolShadowsocks:
				shadowsocks := bo.Shadowsocks{}
				shadowsocks.Name = item.Name
				shadowsocks.Type = constant.ClashShadowsocks
				shadowsocks.Server = item.Domain
				shadowsocks.Port = item.Port
				shadowsocks.Cipher = item.XraySSMethod
				shadowsocks.Password = pass
				shadowsocks.Udp = true
				ClashConfigInterface = append(ClashConfigInterface, shadowsocks)
				proxies = append(proxies, item.Name)
			}
		} else if item.NodeTypeId == constant.TrojanGo {
			nodeTrojanGo, err := SelectNodeTrojanGoById(&item.NodeSubId)
			if err != nil {
				return nil, "", []byte{}, vo.SystemVo{}, err
			}
			trojanGo := bo.TrojanGo{}
			trojanGo.Name = item.Name
			trojanGo.Type = constant.ClashTrojan
			trojanGo.Server = item.Domain
			trojanGo.Port = item.Port
			trojanGo.Password = pass
			trojanGo.Udp = true
			trojanGo.SNI = *nodeTrojanGo.Sni
			if *nodeTrojanGo.WebsocketEnable == 1 {
				trojanGo.Network = "ws"
				trojanGo.WsOpts.Path = *nodeTrojanGo.WebsocketPath
				trojanGo.WsOpts.Headers.Host = *nodeTrojanGo.WebsocketHost
			}
			ClashConfigInterface = append(ClashConfigInterface, trojanGo)
			proxies = append(proxies, item.Name)
		} else if item.NodeTypeId == constant.Hysteria {
			nodeHysteria, err := SelectNodeHysteriaById(&item.NodeSubId)
			if err != nil {
				return nil, "", []byte{}, vo.SystemVo{}, err
			}
			hysteria := bo.Hysteria{}
			hysteria.Name = item.Name
			hysteria.Type = constant.ClashSHysteria
			hysteria.Server = item.Domain
			hysteria.Port = item.Port
			hysteria.AuthStr = pass
			hysteria.Obfs = *nodeHysteria.Obfs
			hysteria.Protocol = *nodeHysteria.Protocol
			hysteria.Up = *nodeHysteria.UpMbps
			hysteria.Down = *nodeHysteria.DownMbps
			ClashConfigInterface = append(ClashConfigInterface, hysteria)
			proxies = append(proxies, item.Name)
		}
	}
	proxyGroups := make([]bo.ProxyGroup, 0)
	proxyGroup := bo.ProxyGroup{
		Name:    "PROXY",
		Type:    "select",
		Proxies: proxies,
	}
	proxyGroups = append(proxyGroups, proxyGroup)
	clashConfig.ProxyGroups = proxyGroups
	clashConfig.Proxies = ClashConfigInterface

	clashConfigYaml, err := yaml.Marshal(&clashConfig)
	if err != nil {
		return nil, "", []byte{}, vo.SystemVo{}, errors.New(constant.SysError)
	}

	systemName := constant.SystemName
	systemConfig, err := SelectSystemByName(&systemName)
	if err != nil {
		return nil, "", []byte{}, vo.SystemVo{}, errors.New(constant.SysError)
	}
	return account, userInfo, clashConfigYaml, systemConfig, nil
}

func ExportAccount(accountId uint, accountUsername string) error {
	fileName := fmt.Sprintf("accountExport-%s.json", time.Now().Format("20060102150405"))
	filePath := fmt.Sprintf("%s/%s", constant.ExportPath, fileName)

	var fileTaskType uint = constant.TaskTypeAccountExport
	var fileTaskStatus = constant.TaskDoing
	fileTask := module.FileTask{
		Name:            &fileName,
		Path:            &filePath,
		Type:            &fileTaskType,
		Status:          &fileTaskStatus,
		AccountId:       &accountId,
		AccountUsername: &accountUsername,
	}
	fileTaskId, err := dao.CreateFileTask(&fileTask)
	if err != nil {
		return err
	}

	go func() {
		var mutex sync.Mutex
		defer mutex.Unlock()
		if mutex.TryLock() {
			var fail = constant.TaskFail
			var success = constant.TaskSuccess
			fileTask := module.FileTask{
				Id:     &fileTaskId,
				Status: &fail,
			}

			// 查询所有需要导出数据
			accountExportVo, err := dao.SelectAccountAll()
			if err != nil {
				logrus.Errorf("ExportAccount SelectAccountAll err: %v", err)
			}
			if err = util.ExportJson(filePath, accountExportVo); err != nil {
				logrus.Errorf("ExportAccount ExportJson err: %v", err)
			} else {
				fileTask.Status = &success
			}

			// 更新文件任务状态
			if err = dao.UpdateFileTaskById(&fileTask); err != nil {
				logrus.Errorf("ExportAccount UpdateFileTaskById err: %v", err)
			}
		}
	}()

	return nil
}

func ImportAccount(cover uint, file *multipart.FileHeader, accountId uint, accountUsername string) error {
	fileName := file.Filename

	var fileTaskType uint = constant.TaskTypeAccountImport
	var fileTaskStatus = constant.TaskDoing
	fileTask := module.FileTask{
		Name:            &fileName,
		Path:            nil,
		Type:            &fileTaskType,
		Status:          &fileTaskStatus,
		AccountId:       &accountId,
		AccountUsername: &accountUsername,
	}
	fileTaskId, err := dao.CreateFileTask(&fileTask)
	if err != nil {
		return err
	}

	go func(fileTaskId uint) {
		var fail = constant.TaskFail
		var success = constant.TaskSuccess
		fileTask := module.FileTask{
			Id:     &fileTaskId,
			Status: &fail,
		}

		src, err := file.Open()
		defer src.Close()
		if err != nil {
			logrus.Errorf("ImportAccount file Open err: %v", err)
			fileUploadError := constant.FileUploadError
			fileTask.ErrMsg = &fileUploadError
			if err = dao.UpdateFileTaskById(&fileTask); err != nil {
				logrus.Errorf("ImportAccount UpdateFileTaskById err: %v", err)
			}
			return
		}

		var accounts []module.Account
		decoder := json.NewDecoder(src)
		if err = decoder.Decode(&accounts); err != nil {
			logrus.Errorf("ImportAccount decoder Decode err: %v", err)
			fileUploadError := constant.FileUploadError
			fileTask.ErrMsg = &fileUploadError
			if err = dao.UpdateFileTaskById(&fileTask); err != nil {
				logrus.Errorf("ImportAccount UpdateFileTaskById err: %v", err)
			}
			return
		}
		if len(accounts) == 0 {
			logrus.Errorf("ImportAccount err: %s", constant.RowNotEnough)
			fileUploadError := constant.RowNotEnough
			fileTask.ErrMsg = &fileUploadError
			if err = dao.UpdateFileTaskById(&fileTask); err != nil {
				logrus.Errorf("ImportAccount UpdateFileTaskById err: %v", err)
			}
			return
		}
		// 在这里可以处理数据并将其存储到数据库中 todo 这里可能存在性能问题
		for _, item := range accounts {
			if err = dao.CreateOrUpdateAccount(item, cover); err != nil {
				continue
			}
		}

		fileTask.Status = &success
		// 更新文件任务状态
		if err = dao.UpdateFileTaskById(&fileTask); err != nil {
			logrus.Errorf("ImportAccount UpdateFileTaskById err: %v", err)
		}
	}(fileTaskId)
	return nil
}

func LoginLimit(username string) {
	redis.Client.String.
		Incr(fmt.Sprintf("trojan-panel:login-limit:%s", username))
}

// LoginVerify 密码输入错误3次以上 将账户锁定30分钟
func LoginVerify(username string) error {
	get := redis.Client.String.
		Get(fmt.Sprintf("trojan-panel:login-limit:%s", username))
	reply, err := get.Result()
	if err != nil {
		return errors.New(constant.SysError)
	}
	if reply != nil {
		result, err := get.Int()
		if err != nil {
			return errors.New(constant.SysError)
		}
		if result >= 3 {
			redis.Client.String.Set(fmt.Sprintf("trojan-panel:login-limit:%s", username), -1, time.Minute.Milliseconds()*30/1000)
		}
		if result >= 3 || result == -1 {
			return errors.New(constant.LoginLimitError)
		}
	}
	return nil
}
