package api

import (
	"encoding/base64"
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
	var accountLoginDto dto.AccountLoginDto
	_ = c.ShouldBindJSON(&accountLoginDto)
	if err := validate.Struct(&accountLoginDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	account, err := service.SelectAccountByUsernameAndPass(accountLoginDto.Username, accountLoginDto.Pass)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	if *account.Deleted == 1 {
		vo.Fail(constant.AccountDisabled, c)
		return
	}
	if account != nil {
		accountVo := vo.AccountVo{
			Id:       *account.Id,
			Username: *account.Username,
			RoleId:   *account.RoleId,
			Deleted:  *account.Deleted,
		}
		tokenStr, err := util.GenToken(accountVo)
		if err != nil {
			vo.Fail(constant.SysError, c)
		} else {
			if _, err := redis.Client.String.
				Set(fmt.Sprintf("trojan-panel:token:%s", *accountLoginDto.Username), tokenStr,
					time.Hour.Milliseconds()*2/1000).Result(); err != nil {
				vo.Fail(constant.SysError, c)
			} else {
				accountLoginVo := vo.AccountLoginVo{
					Token: tokenStr,
				}
				vo.Success(accountLoginVo, c)
			}
		}
		return
	}
	vo.Fail(constant.UsernameOrPassError, c)
}

// HysteriaApi hysteria api
func HysteriaApi(c *gin.Context) {
	var hysteriaAutoDto dto.HysteriaAutoDto
	_ = c.ShouldBindJSON(&hysteriaAutoDto)
	if err := validate.Struct(&hysteriaAutoDto); err != nil {
		vo.HysteriaApiFail(constant.ValidateFailed, c)
		return
	}
	decodeString, err := base64.StdEncoding.DecodeString(*hysteriaAutoDto.Payload)
	if err != nil {
		vo.HysteriaApiFail(constant.ValidateFailed, c)
		return
	}
	usernameAndPass := strings.Split(string(decodeString), "&")
	account, err := service.SelectAccountByUsernameAndPass(&usernameAndPass[0], &usernameAndPass[1])
	if err != nil {
		vo.HysteriaApiFail(err.Error(), c)
		return
	}
	if *account.Deleted == 1 {
		vo.HysteriaApiFail(constant.AccountDisabled, c)
		return
	}
	if account != nil {
		vo.HysteriaApiSuccess(*account.Username, c)
		return
	}
	vo.HysteriaApiFail(constant.UsernameOrPassError, c)
}

// GenerateCaptcha 验证码
func GenerateCaptcha(c *gin.Context) {
	return
}

func Register(c *gin.Context) {
	var accountRegisterDto dto.AccountRegisterDto
	_ = c.ShouldBindJSON(&accountRegisterDto)
	if err := validate.Struct(&accountRegisterDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	if err := service.Register(accountRegisterDto); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func Logout(c *gin.Context) {
	account := util.GetCurrentAccount(c)
	if _, err := redis.Client.Key.
		Del(fmt.Sprintf("trojan-panel:token:%s", account.Username)).
		Result(); err != nil {
		vo.Fail(constant.LogOutError, c)
		return
	}
	vo.Success(nil, c)
}

func CreateAccount(c *gin.Context) {
	var accountCreateDto dto.AccountCreateDto
	_ = c.ShouldBindJSON(&accountCreateDto)
	if err := validate.Struct(&accountCreateDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	if err := service.CreateAccount(accountCreateDto); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func SelectAccountById(c *gin.Context) {
	var accountRequiredIdDto dto.RequiredIdDto
	_ = c.ShouldBindQuery(&accountRequiredIdDto)
	if err := validate.Struct(&accountRequiredIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	account, err := service.SelectAccountById(accountRequiredIdDto.Id)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	accountVo := vo.AccountVo{
		Id:         *account.Id,
		Username:   *account.Username,
		RoleId:     *account.RoleId,
		Email:      *account.Email,
		ExpireTime: *account.ExpireTime,
		Deleted:    *account.Deleted,
		Quota:      *account.Quota,
		Download:   *account.Download,
		Upload:     *account.Upload,
	}
	vo.Success(accountVo, c)
}

func SelectAccountPage(c *gin.Context) {
	var accountPageDto dto.AccountPageDto
	_ = c.ShouldBindQuery(&accountPageDto)
	if err := validate.Struct(&accountPageDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	page, err := service.SelectAccountPage(accountPageDto.Username, accountPageDto.PageNum, accountPageDto.PageSize)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(page, c)
}

func DeleteAccountById(c *gin.Context) {
	var accountRequiredIdDto dto.RequiredIdDto
	_ = c.ShouldBindJSON(&accountRequiredIdDto)
	if err := validate.Struct(&accountRequiredIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	account, err := service.SelectAccountById(accountRequiredIdDto.Id)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	if *account.RoleId == constant.SYSADMIN {
		vo.Fail("不能删除超级管理员账户", c)
		return
	}
	if err := service.DeleteAccountById(accountRequiredIdDto.Id); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func UpdateAccountProfile(c *gin.Context) {
	var accountUpdateProfileDto dto.AccountUpdateProfileDto
	_ = c.ShouldBindJSON(&accountUpdateProfileDto)
	if err := validate.Struct(&accountUpdateProfileDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	if err := service.UpdateAccountProfile(accountUpdateProfileDto.OldPass, accountUpdateProfileDto.NewPass,
		accountUpdateProfileDto.Username, accountUpdateProfileDto.Email); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func GetAccountInfo(c *gin.Context) {
	accountInfo, err := service.GetAccountInfo(c)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(accountInfo, c)
}

func UpdateAccountById(c *gin.Context) {
	var accountUpdateDto dto.AccountUpdateDto
	_ = c.ShouldBindJSON(&accountUpdateDto)
	if err := validate.Struct(&accountUpdateDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	toByte := util.ToByte(*accountUpdateDto.Quota)
	account := module.Account{
		Id:                 accountUpdateDto.Id,
		Quota:              &toByte,
		Username:           accountUpdateDto.Username,
		Pass:               accountUpdateDto.Pass,
		Email:              accountUpdateDto.Email,
		RoleId:             accountUpdateDto.RoleId,
		Deleted:            accountUpdateDto.Deleted,
		ExpireTime:         accountUpdateDto.ExpireTime,
		IpLimit:            accountUpdateDto.IpLimit,
		UploadSpeedLimit:   accountUpdateDto.UploadSpeedLimit,
		DownloadSpeedLimit: accountUpdateDto.DownloadSpeedLimit,
	}
	if err := service.UpdateAccountById(&account); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}
