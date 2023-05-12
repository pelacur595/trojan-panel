package api

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"strconv"
	"strings"
	"time"
	"trojan-panel/dao"
	"trojan-panel/dao/redis"
	"trojan-panel/module"
	"trojan-panel/module/constant"
	"trojan-panel/module/dto"
	"trojan-panel/module/vo"
	"trojan-panel/service"
	"trojan-panel/util"
)

func Login(c *gin.Context) {
	var accountLoginDto dto.AccountLoginDto
	_ = c.ShouldBindJSON(&accountLoginDto)
	if err := validate.Struct(&accountLoginDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	if err := service.LoginVerify(*accountLoginDto.Username); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	systemName := constant.SystemName
	systemVo, err := service.SelectSystemByName(&systemName)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	if systemVo.CaptchaEnable == 1 && !util.VerifyCaptcha(*accountLoginDto.CaptchaId, *accountLoginDto.CaptchaCode) {
		vo.Fail(constant.CaptchaError, c)
		return
	}
	account, err := service.SelectAccountByUsername(accountLoginDto.Username)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	if account != nil {
		if !util.Sha1Match(*account.Pass, fmt.Sprintf("%s%s", *accountLoginDto.Username, *accountLoginDto.Pass)) {
			vo.Fail(constant.UsernameOrPassError, c)
			service.LoginLimit(*account.Username)
			return
		}
		if *account.Deleted != 0 {
			vo.Fail(constant.AccountDisabled, c)
			return
		}
		roles, err := dao.SelectRoleNameByParentId(account.RoleId, true)
		if err != nil {
			vo.Fail(constant.SysError, c)
			return
		}
		accountVo := vo.AccountVo{
			Id:       *account.Id,
			Username: *account.Username,
			RoleId:   *account.RoleId,
			Deleted:  *account.Deleted,
			Roles:    roles,
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
				milli := uint(time.Now().UnixMilli())
				// 第一次登录
				account := module.Account{
					Id:            account.Id,
					LastLoginTime: &milli,
				}
				if account.ValidityPeriod != nil && *account.ValidityPeriod >= 0 &&
					account.LastLoginTime != nil && *account.LastLoginTime == 0 {
					expireTime := milli + *account.ValidityPeriod*24*60*60*1000
					account.ExpireTime = &expireTime
				}
				if err := service.UpdateAccountById(tokenStr, &account); err != nil {
					vo.Fail(constant.SysError, c)
					return
				}
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

// GenerateCaptcha 验证码
func GenerateCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverMath(80, 240, 8, 3, nil, nil, []string{"wqy-microhei.ttc"})
	captcha := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	// 生成验证码图片和答案
	id, b64s, err := captcha.Generate()
	if err != nil {
		vo.Fail(constant.CaptchaGenerateError, c)
		return
	}
	captureVo := vo.CaptureVo{
		CaptchaId:  id,
		CaptchaImg: b64s,
	}
	vo.Success(captureVo, c)
}

func Register(c *gin.Context) {
	var accountRegisterDto dto.AccountRegisterDto
	_ = c.ShouldBindJSON(&accountRegisterDto)
	if err := validate.Struct(&accountRegisterDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	systemName := constant.SystemName
	systemVo, err := service.SelectSystemByName(&systemName)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	if systemVo.CaptchaEnable == 1 && !util.VerifyCaptcha(*accountRegisterDto.CaptchaId, *accountRegisterDto.CaptchaCode) {
		vo.Fail(constant.CaptchaError, c)
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
		Id:             *account.Id,
		Username:       *account.Username,
		RoleId:         *account.RoleId,
		Email:          *account.Email,
		ValidityPeriod: *account.ValidityPeriod,
		ExpireTime:     *account.ExpireTime,
		Deleted:        *account.Deleted,
		Quota:          *account.Quota,
		Download:       *account.Download,
		Upload:         *account.Upload,
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
	page, err := service.SelectAccountPage(accountPageDto.Username, accountPageDto.Deleted, accountPageDto.PageNum,
		accountPageDto.PageSize)
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
	if err := service.DeleteAccountById(util.GetToken(c), accountRequiredIdDto.Id); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func UpdateAccountPass(c *gin.Context) {
	var accountUpdatePassDto dto.AccountUpdatePassDto
	_ = c.ShouldBindJSON(&accountUpdatePassDto)
	if err := validate.Struct(&accountUpdatePassDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	account := util.GetCurrentAccount(c)
	if err := service.UpdateAccountPass(util.GetToken(c), accountUpdatePassDto.OldPass, accountUpdatePassDto.NewPass,
		&account.Username); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func UpdateAccountProperty(c *gin.Context) {
	var accountUpdatePropertyDto dto.AccountUpdatePropertyDto
	_ = c.ShouldBindJSON(&accountUpdatePropertyDto)
	if err := validate.Struct(&accountUpdatePropertyDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	account := util.GetCurrentAccount(c)
	if err := service.UpdateAccountProperty(util.GetToken(c), &account.Username,
		accountUpdatePropertyDto.Pass, accountUpdatePropertyDto.Username, accountUpdatePropertyDto.Email); err != nil {
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
		Id:             accountUpdateDto.Id,
		Quota:          &toByte,
		Username:       accountUpdateDto.Username,
		Pass:           accountUpdateDto.Pass,
		Email:          accountUpdateDto.Email,
		RoleId:         accountUpdateDto.RoleId,
		Deleted:        accountUpdateDto.Deleted,
		ValidityPeriod: accountUpdateDto.ValidityPeriod,
		ExpireTime:     accountUpdateDto.ExpireTime,
		//IpLimit:            accountUpdateDto.IpLimit,
		//UploadSpeedLimit:   accountUpdateDto.UploadSpeedLimit,
		//DownloadSpeedLimit: accountUpdateDto.DownloadSpeedLimit,
	}
	if err := service.UpdateAccountById(util.GetToken(c), &account); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func ClashSubscribe(c *gin.Context) {
	accountVo := util.GetCurrentAccount(c)
	password, err := service.SelectConnectPassword(&accountVo.Id, &accountVo.Username)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(fmt.Sprintf("/api/auth/subscribe/%s", base64.StdEncoding.EncodeToString([]byte(password))), c)
}

// Subscribe 订阅
func Subscribe(c *gin.Context) {
	token := c.Param("token")
	//userAgent := c.Request.Header.Get("User-Agent")
	tokenDecode, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	pass := string(tokenDecode)

	//if strings.HasPrefix(userAgent, constant.ClashforWindows) {
	account, userInfo, clashConfigYaml, systemConfig, err := service.SubscribeClash(pass)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	result := fmt.Sprintf(`%s
%s`, string(clashConfigYaml), systemConfig.ClashRule)

	c.Header("content-disposition", fmt.Sprintf("attachment; filename=%s.yaml", *account.Username))
	c.Header("profile-update-interval", "12")
	c.Header("subscription-userinfo", userInfo)
	c.String(200, result)
	return
	//}
	//vo.Fail("不支持该客户端", c)
}

// ResetAccountDownloadAndUpload 重设下载和上传流量
func ResetAccountDownloadAndUpload(c *gin.Context) {
	var accountRequiredIdDto dto.RequiredIdDto
	_ = c.ShouldBindJSON(&accountRequiredIdDto)
	if err := validate.Struct(&accountRequiredIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	if err := service.ResetAccountDownloadAndUpload(accountRequiredIdDto.Id, nil); err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	vo.Success(nil, c)
}

// ExportAccount 导出用户
func ExportAccount(c *gin.Context) {
	accountVo := util.GetCurrentAccount(c)
	if err := service.ExportAccount(accountVo.Id, accountVo.Username); err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	vo.Success(nil, c)
}

// ImportAccount 导入用户
func ImportAccount(c *gin.Context) {
	coverStr, b := c.GetPostForm("cover")
	if !b {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	cover, err := strconv.ParseUint(coverStr, 10, 32)
	file, err := c.FormFile("file")
	if err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	// 文件大小 10MB
	if file.Size > 1024*1024*10 {
		vo.Fail(constant.FileSizeTooBig, c)
		return
	}
	// 文件后缀.json
	if !strings.HasSuffix(file.Filename, ".json") {
		vo.Fail(constant.FileFormatError, c)
		return
	}
	account := util.GetCurrentAccount(c)
	if err := service.ImportAccount(uint(cover), file, account.Id, account.Username); err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	vo.Success(nil, c)
}

func CreateAccountBatch(c *gin.Context) {
	var createAccountBatchDto dto.CreateAccountBatchDto
	_ = c.ShouldBindJSON(&createAccountBatchDto)
	if err := validate.Struct(&createAccountBatchDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	account := util.GetCurrentAccount(c)
	if err := service.CreateAccountBatch(account.Id, account.Username, createAccountBatchDto); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func ExportAccountUnused(c *gin.Context) {
	accountVo := util.GetCurrentAccount(c)
	if err := service.ExportAccountUnused(accountVo.Id, accountVo.Username); err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	vo.Success(nil, c)
}
