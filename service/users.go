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

func CreateUser(userCreateDto dto.UserCreateDto) error {
	count, err := dao.CountUserByUsername(userCreateDto.Username)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(constant.UsernameExist)
	}
	toByte := util.ToByte(*userCreateDto.Quota)
	user := module.Users{
		Quota:      &toByte,
		Username:   userCreateDto.Username,
		Pass:       userCreateDto.Pass,
		RoleId:     userCreateDto.RoleId,
		Deleted:    userCreateDto.Deleted,
		ExpireTime: userCreateDto.ExpireTime,
		Email:      userCreateDto.Email,
	}
	if err := dao.CreateUser(&user); err != nil {
		return err
	}
	if user.Deleted != nil && *user.Deleted == 1 {
		if err := PullUserWhiteOrBlackByUsername([]string{*user.Username}, true); err != nil {
			return err
		}
	} else if *user.ExpireTime <= util.NowMilli() {
		if err := DisableUsers([]string{*user.Username}); err != nil {
			return err
		}
	}

	return nil
}
func SelectUserById(id *uint) (*vo.UsersVo, error) {
	userVo, err := dao.SelectUserById(id)
	if err != nil {
		return nil, err
	}
	return userVo, nil
}
func CountUserByUsername(username *string) (int, error) {
	count, err := dao.CountUserByUsername(username)
	if err != nil {
		return 0, err
	}
	return count, err
}
func SelectUserPage(username *string, pageNum *uint, pageSize *uint) (*vo.UsersPageVo, error) {
	page, err := dao.SelectUserPage(username, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	return page, err
}
func DeleteUserById(id *uint) error {
	if err := dao.DeleteUserById(id); err != nil {
		return err
	}
	return nil
}

func SelectUserByUsernameAndPass(username *string, pass *string) (*vo.UsersVo, error) {
	userVo, err := dao.SelectUserByUsernameAndPass(username, pass)
	if err != nil {
		return nil, err
	}
	return userVo, nil
}

func UpdateUserPassByUsername(oldPass *string, newPass *string, username *string) error {
	if err := dao.UpdateUserPassByUsername(oldPass, newPass, username); err != nil {
		return err
	}
	return nil
}

// 获取当前请求账户信息
func GetUserInfo(c *gin.Context) (*vo.UserInfo, error) {
	userVo := util.GetCurrentUser(c)
	menuList, err := SelectMenuListByRoleId(&userVo.RoleId)
	if err != nil {
		return nil, err
	}
	roleNames, err := dao.SelectRoleNameByParentId(&userVo.RoleId, true)
	if err != nil {
		return nil, err
	}
	userInfo := vo.UserInfo{
		Id:        userVo.Id,
		Username:  userVo.Username,
		RoleNames: roleNames,
		MenuList:  menuList,
	}
	return &userInfo, nil
}

func UpdateUserById(users *module.Users) error {
	if err := dao.UpdateUserById(users); err != nil {
		return err
	}
	if users.Deleted != nil && *users.Deleted == 1 {
		if err := PullUserWhiteOrBlackByUsername([]string{*users.Username}, true); err != nil {
			return err
		}
	} else if *users.ExpireTime <= util.NowMilli() {
		if err := DisableUsers([]string{*users.Username}); err != nil {
			return err
		}
	}
	return nil
}

func Register(userRegisterDto dto.UserRegisterDto) error {
	name := constant.SystemName
	systemVo, err := SelectSystemByName(&name)
	if err != nil {
		return err
	}
	if systemVo.OpenRegister == 0 {
		return errors.New(constant.UserRegisterClosed)
	}

	count, err := dao.CountUserByUsername(userRegisterDto.Username)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(constant.UsernameExist)
	}
	u := constant.USER
	milli := util.DayToMilli(systemVo.RegisterExpireDays)
	registerQuota := util.ToByte(systemVo.RegisterQuota)
	user := module.Users{
		Quota:      &registerQuota,
		Username:   userRegisterDto.Username,
		Pass:       userRegisterDto.Pass,
		Email:      userRegisterDto.Email,
		RoleId:     &u,
		Deleted:    new(uint),
		ExpireTime: &milli,
	}
	if err := dao.CreateUser(&user); err != nil {
		return err
	}
	return nil
}

// 拉白或者拉黑用户 此操作会清空用户流量
func PullUserWhiteOrBlackByUsername(usernames []string, isBlack bool) error {
	if len(usernames) > 0 {
		var deleted uint
		if isBlack {
			deleted = 1
		} else {
			deleted = 0
		}
		if err := dao.UpdateUserQuotaOrDownloadOrUploadOrDeletedByUsernames(usernames, new(int), new(uint), new(uint), &deleted); err != nil {
			return err
		}
	}
	return nil
}

// 清空流量/禁用用户连接节点
func DisableUsers(usernames []string) error {
	if len(usernames) > 0 {
		if err := dao.UpdateUserQuotaOrDownloadOrUploadOrDeletedByUsernames(usernames, new(int), new(uint), new(uint), nil); err != nil {
			return err
		}
	}
	return nil
}

// 定时任务：扫描无效用户
func ScanUsers() {
	usernames, err := dao.SelectUsernameByDeletedOrExpireTime()
	if err != nil {
		logrus.Errorln(err.Error())
		return
	}

	if len(usernames) > 0 {
		if err := DisableUsers(usernames); err != nil {
			logrus.Errorf("定时扫描用户任务禁用用户异常 usernames: %s error: %v\n", usernames, err)
		}
		logrus.Infof("定时扫描用户任务禁用用户 usernames: %s\n", usernames)
	}
}

// 定时任务：到期警告
func ScanUserExpireWarn() {
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
	users, err := dao.SelectUsersEmailByExpireTime(util.DayToMilli(expireWarnDay))
	if err != nil {
		logrus.Errorln(err.Error())
		return
	}
	if len(users) > 0 {
		for _, user := range users {
			if user.Email != nil && *user.Email != "" {
				// 发送到期邮件
				emailDto := dto.SendEmailDto{
					FromEmailName: "Trojan Panel",
					ToEmails:      []string{*user.Email},
					Subject:       "账号到期提醒",
					Content:       fmt.Sprintf("您的账户: %s,还有%d天到期,请及时续期", *user.Username, expireWarnDay),
				}
				if err := SendEmail(&emailDto); err != nil {
					logrus.Errorln(fmt.Sprintf("到期警告邮件发送失败 err: %v", err))
				}
			}
		}
	}
}
