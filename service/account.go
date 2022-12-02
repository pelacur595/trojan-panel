package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"sync"
	"trojan-panel/core"
	"trojan-panel/dao"
	"trojan-panel/module"
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
func SelectAccountPage(username *string, pageNum *uint, pageSize *uint) (*vo.AccountPageVo, error) {
	return dao.SelectAccountPage(username, pageNum, pageSize)
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

func UpdateAccountProfile(token string, oldPass *string, newPass *string, username *string, email *string) error {
	var mutex sync.Mutex
	defer mutex.Unlock()
	if mutex.TryLock() {
		if oldPass != nil && *oldPass != "" && newPass != nil && *newPass != "" {
			password, err := dao.SelectConnectPassword(nil, username)
			if err != nil {
				return err
			}
			if err = RemoveAccount(token, password); err != nil {
				return err
			}
		}
		if err := dao.UpdateAccountProfile(oldPass, newPass, username, email); err != nil {
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
					FromEmailName: "Trojan Panel",
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

// ClashSubscribe 导出Clash配置
func ClashSubscribe(pass string) (*module.Account, []vo.NodeOneVo, error) {
	account, err := dao.SelectAccountClashSubscribe(pass)
	if err != nil {
		return nil, nil, err
	}
	nodes, err := dao.SelectNodesIpAndPort()
	var nodeVos []vo.NodeOneVo
	for _, item := range nodes {
		nodeOneVo, err := SelectNodeById(item.Id)
		if err != nil {
			continue
		}
		nodeVos = append(nodeVos, *nodeOneVo)
	}
	return account, nodeVos, nil
}

func RemoveAccount(token string, password string) error {
	ips, err := dao.SelectNodesIpDistinct()
	if err != nil {
		return err
	}
	for _, ip := range ips {
		removeDto := core.AccountRemoveDto{
			Password: password,
		}
		_ = core.RemoveAccount(ip, token, &removeDto)
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
