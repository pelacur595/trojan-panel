package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
	"trojan/dao/redis"
	"trojan/module"
	"trojan/module/constant"
	"trojan/module/dto"
	"trojan/module/vo"
	"trojan/service"
	"trojan/util"
)

func Login(c *gin.Context) {
	var userLoginDto dto.UserLoginDto
	_ = c.ShouldBindJSON(&userLoginDto)
	if err := validate.Struct(&userLoginDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	usersVo, err := service.SelectUserByUsernameAndPass(userLoginDto.Username, userLoginDto.Pass)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	if usersVo.Deleted == 1 {
		vo.Fail(constant.UserDisabled, c)
		return
	}
	if usersVo != nil {
		tokenStr, err := util.GenToken(*usersVo)
		if err != nil {
			vo.Fail(constant.SysError, c)
		} else {
			if _, err := redis.Client.String.
				Set(fmt.Sprintf("trojan-panel:token:%s", *userLoginDto.Username), tokenStr,
					time.Hour.Milliseconds()*2/1000).Result(); err != nil {
				vo.Fail(constant.SysError, c)
			} else {
				userLoginVo := vo.UsersLoginVo{
					Token: tokenStr,
				}
				vo.Success(userLoginVo, c)
			}
		}
		return
	}
	vo.Fail(constant.UsernameOrPassError, c)
}

// hysteria api
func HysteriaApi(c *gin.Context) {
	var hysteriaAutoDto dto.HysteriaAutoDto
	_ = c.ShouldBindJSON(&hysteriaAutoDto)
	if err := validate.Struct(&hysteriaAutoDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	usernameAndPass := strings.Split(*hysteriaAutoDto.Payload, "&")
	usersVo, err := service.SelectUserByUsernameAndPass(&usernameAndPass[0], &usernameAndPass[1])
	if err != nil {
		vo.HysteriaApiFail(err.Error(), c)
		return
	}
	if usersVo.Deleted == 1 {
		vo.HysteriaApiFail(constant.UserDisabled, c)
		return
	}
	if usersVo != nil {
		vo.HysteriaApiSuccess(usersVo.Username, c)
		return
	}
	vo.HysteriaApiFail(constant.UsernameOrPassError, c)
}

// 验证码
func GenerateCaptcha(c *gin.Context) {
	return
}

func Register(c *gin.Context) {
	var userRegisterDto dto.UserRegisterDto
	_ = c.ShouldBindJSON(&userRegisterDto)
	if err := validate.Struct(&userRegisterDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	if err := service.Register(userRegisterDto); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func Logout(c *gin.Context) {
	user := util.GetCurrentUser(c)
	if _, err := redis.Client.Key.
		Del(fmt.Sprintf("trojan-panel:token:%s", user.Username)).
		Result(); err != nil {
		vo.Fail(constant.LogOutError, c)
		return
	}
	vo.Success(nil, c)
}

func CreateUser(c *gin.Context) {
	var userCreateDto dto.UserCreateDto
	_ = c.ShouldBindJSON(&userCreateDto)
	if err := validate.Struct(&userCreateDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	if err := service.CreateUser(userCreateDto); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func SelectUserById(c *gin.Context) {
	var userRequiredIdDto dto.RequiredIdDto
	_ = c.ShouldBindQuery(&userRequiredIdDto)
	if err := validate.Struct(&userRequiredIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	userVo, err := service.SelectUserById(userRequiredIdDto.Id)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(userVo, c)
}

func SelectUserPage(c *gin.Context) {
	var userPageDto dto.UsersPageDto
	_ = c.ShouldBindQuery(&userPageDto)
	if err := validate.Struct(&userPageDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	page, err := service.SelectUserPage(userPageDto.Username, userPageDto.PageNum, userPageDto.PageSize)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(page, c)
}

func DeleteUserById(c *gin.Context) {
	var userRequiredIdDto dto.RequiredIdDto
	_ = c.ShouldBindJSON(&userRequiredIdDto)
	if err := validate.Struct(&userRequiredIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	userVo, err := service.SelectUserById(userRequiredIdDto.Id)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	if userVo.RoleId == constant.SYSADMIN {
		vo.Fail("不能删除超级管理员账户", c)
		return
	}
	if err := service.DeleteUserById(userRequiredIdDto.Id); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func UpdateUserProfile(c *gin.Context) {
	var userUpdateProfileDto dto.UserUpdateProfileDto
	_ = c.ShouldBindJSON(&userUpdateProfileDto)
	if err := validate.Struct(&userUpdateProfileDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	if err := service.UpdateUserProfile(userUpdateProfileDto.OldPass, userUpdateProfileDto.NewPass,
		userUpdateProfileDto.Username, userUpdateProfileDto.Email); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func GetUserInfo(c *gin.Context) {
	user, err := service.GetUserInfo(c)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(user, c)
}

func UpdateUserById(c *gin.Context) {
	var userUpdateDto dto.UserUpdateDto
	_ = c.ShouldBindJSON(&userUpdateDto)
	if err := validate.Struct(&userUpdateDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	toByte := util.ToByte(*userUpdateDto.Quota)
	users := module.Users{
		Id:         userUpdateDto.Id,
		Quota:      &toByte,
		Username:   userUpdateDto.Username,
		Pass:       userUpdateDto.Pass,
		Email:      userUpdateDto.Email,
		RoleId:     userUpdateDto.RoleId,
		Deleted:    userUpdateDto.Deleted,
		ExpireTime: userUpdateDto.ExpireTime,
		//IpLimit:            userUpdateDto.IpLimit,
		//UploadSpeedLimit:   userUpdateDto.UploadSpeedLimit,
		//DownloadSpeedLimit: userUpdateDto.DownloadSpeedLimit,
	}
	if err := service.UpdateUserById(&users); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}
