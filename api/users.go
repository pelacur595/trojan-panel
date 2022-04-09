package api

import (
	"github.com/gin-gonic/gin"
	"log"
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
			userLoginVo := vo.UsersLoginVo{
				Token: tokenStr,
			}
			vo.Success(userLoginVo, c)
		}
		return
	}
	vo.Fail(constant.UsernameOrPassError, c)
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

func CreateUser(c *gin.Context) {
	var userCreateDto dto.UserCreateDto
	_ = c.ShouldBindJSON(&userCreateDto)
	if err := validate.Struct(&userCreateDto); err != nil {
		log.Println(err.Error())
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
		log.Println(err)
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

func UpdateUserPassByUsername(c *gin.Context) {
	var userUpdatePassDto dto.UserUpdatePassDto
	_ = c.ShouldBindJSON(&userUpdatePassDto)
	if err := validate.Struct(&userUpdatePassDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	if err := service.UpdateUserPassByUsername(userUpdatePassDto.OldPass, userUpdatePassDto.NewPass, userUpdatePassDto.Username); err != nil {
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
		log.Println(err)
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	toByte := util.ToByte(*userUpdateDto.Quota)
	users := module.Users{
		Id:         userUpdateDto.Id,
		Quota:      &toByte,
		Username:   userUpdateDto.Username,
		Pass:       userUpdateDto.Pass,
		RoleId:     userUpdateDto.RoleId,
		Deleted:    userUpdateDto.Deleted,
		ExpireTime: userUpdateDto.ExpireTime,
	}
	if err := service.UpdateUserById(&users); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}
