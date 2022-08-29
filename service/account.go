package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"trojan/dao"
	"trojan/module"
	"trojan/module/constant"
	"trojan/module/dto"
	"trojan/module/vo"
	"trojan/util"
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
	}
	if err := dao.CreateAccount(&account); err != nil {
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
func DeleteAccountById(id *uint) error {
	return dao.DeleteAccountById(id)
}

func SelectAccountByUsernameAndPass(username *string, pass *string) (*module.Account, error) {
	return dao.SelectAccountByUsernameAndPass(username, pass)
}

func UpdateAccountProfile(oldPass *string, newPass *string, username *string, email *string) error {
	return dao.UpdateAccountProfile(oldPass, newPass, username, email)
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

func UpdateAccountById(account *module.Account) error {
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
	return nil
}

func Register(accountRegisterDto dto.AccountRegisterDto) error {
	name := constant.SystemName
	systemVo, err := SelectSystemByName(&name)
	if err != nil {
		return err
	}
	if systemVo.OpenRegister == 0 {
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
	if err := dao.CreateAccount(&account); err != nil {
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

// ScanAccounts 定时任务：扫描无效用户
func ScanAccounts() {
	usernames, err := dao.SelectAccountUsernameByDeletedOrExpireTime()
	if err != nil {
		logrus.Errorln(err.Error())
		return
	}

	if len(usernames) > 0 {
		if err := DisableAccount(usernames); err != nil {
			logrus.Errorf("定时扫描用户任务禁用用户异常 usernames: %s error: %v\n", usernames, err)
		}
		logrus.Infof("定时扫描用户任务禁用用户 usernames: %s\n", usernames)
	}
}

// ScanAccountResetTraffic 设置角色为普通用户流量为0
func ScanAccountResetTraffic() {
	if err := dao.UpdateAccountQuota(); err != nil {
		logrus.Errorf("定时重设用户总流量异常 error: %v\n", err)
	}
}

// ScanAccountExpireWarn 定时任务：到期警告
func ScanAccountExpireWarn() {
	systemName := constant.SystemName
	systemVo, err := SelectSystemByName(&systemName)
	if err != nil {
		logrus.Errorln(err.Error())
		return
	}
	if systemVo.EmailEnable == 0 || systemVo.ExpireWarnEnable == 0 {
		return
	}
	expireWarnDay := systemVo.ExpireWarnDay
	accounts, err := dao.SelectAccountsEmailByExpireTime(util.DayToMilli(expireWarnDay))
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
				if err := SendEmail(&emailDto); err != nil {
					logrus.Errorln(fmt.Sprintf("到期警告邮件发送失败 err: %v", err))
				}
			}
		}
	}
}
