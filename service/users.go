package service

import (
	"errors"
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
		Email:      userCreateDto.Email,
		RoleId:     userCreateDto.RoleId,
		Deleted:    userCreateDto.Deleted,
		ExpireTime: userCreateDto.ExpireTime,
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
	} else {
		if err := EnableUser(users.Username); err != nil {
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
	var deleted uint = 0
	milli := util.DayToMilli(systemVo.RegisterExpireDays)
	registerQuota := util.ToByte(systemVo.RegisterQuota)
	user := module.Users{
		Quota:      &registerQuota,
		Username:   userRegisterDto.Username,
		Pass:       userRegisterDto.Pass,
		Email:      userRegisterDto.Email,
		RoleId:     &u,
		Deleted:    &deleted,
		ExpireTime: &milli,
	}
	if err := dao.CreateUser(&user); err != nil {
		return err
	}
	return nil
}

// 拉白或者拉黑用户
func PullUserWhiteOrBlackByUsername(usernames []string, isBlack bool) error {
	if len(usernames) > 0 {
		if isBlack {
			// 拉黑
			password := ""
			var black uint = 1
			if err := dao.UpdateUserPasswordOrDeletedByUsernames(usernames, &password, &black); err != nil {
				return err
			}
		} else {
			// 拉白
			var white uint = 0
			if err := dao.UpdateUserPasswordOrDeletedByUsernames(usernames, nil, &white); err != nil {
				return err
			}
			for _, username := range usernames {
				if err := EnableUser(&username); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// 启用用户服务链接
func EnableUser(username *string) error {
	if username != nil && *username != "" {
		encryPassword, err := dao.SelectEncryPasswordByUsername(username)
		if err != nil {
			return err
		}
		if err := dao.UpdateUserPasswordOrDeletedByUsernames([]string{*username}, &encryPassword, nil); err != nil {
			return err
		}
	}
	return nil
}

// 禁用用户服务链接
func DisableUsers(usernames []string) error {
	if len(usernames) > 0 {
		password := ""
		if err := dao.UpdateUserPasswordOrDeletedByUsernames(usernames, &password, nil); err != nil {
			return err
		}
	}
	return nil
}

// 扫描无效用户
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
