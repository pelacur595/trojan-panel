package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"time"
	"trojan-panel/dao"
	"trojan-panel/dao/redis"
	"trojan-panel/module"
	"trojan-panel/module/bo"
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
	account, err := service.SelectAccountByUsername(accountLoginDto.Username)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	if !util.Sha1Match(*account.Pass, fmt.Sprintf("%s%s", *accountLoginDto.Username, *accountLoginDto.Pass)) {
		vo.Fail(constant.UsernameOrPassError, c)
		return
	}
	if *account.Deleted != 0 {
		vo.Fail(constant.AccountDisabled, c)
		return
	}
	if account != nil {
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
	if err := service.DeleteAccountById(util.GetToken(c), accountRequiredIdDto.Id); err != nil {
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
	if err := service.UpdateAccountProfile(util.GetToken(c), accountUpdateProfileDto.OldPass, accountUpdateProfileDto.NewPass,
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
		Id:         accountUpdateDto.Id,
		Quota:      &toByte,
		Username:   accountUpdateDto.Username,
		Pass:       accountUpdateDto.Pass,
		Email:      accountUpdateDto.Email,
		RoleId:     accountUpdateDto.RoleId,
		Deleted:    accountUpdateDto.Deleted,
		ExpireTime: accountUpdateDto.ExpireTime,
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
	vo.Success(fmt.Sprintf("/api/auth/clash/%s", base64.StdEncoding.EncodeToString([]byte(password))), c)
}

// Clash
// Clash for windows 参考文档：
// 1. https://docs.cfw.lbyczf.com/contents/urlscheme.html
// 2. https://github.com/crossutility/Quantumult/blob/master/extra-subscription-feature.md
// 3. https://github.com/Dreamacro/clash/wiki/Configuration
func Clash(c *gin.Context) {
	token := c.Param("token")
	tokenDecode, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	pass := string(tokenDecode)
	account, nodeOneVos, err := service.ClashSubscribe(pass)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
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
		if item.NodeTypeId == 1 {
			nodeXray, err := service.SelectNodeXrayById(&item.NodeSubId)
			if err != nil {
				vo.Fail(err.Error(), c)
				return
			}

			streamSettings := bo.StreamSettings{}
			if nodeXray.StreamSettings != nil && *nodeXray.StreamSettings != "" {
				if err = json.Unmarshal([]byte(*nodeXray.StreamSettings), &streamSettings); err != nil {
					logrus.Errorln(fmt.Sprintf("SystemVo JSON反转失败 err: %v", err))
					vo.Fail(constant.SysError, c)
					return
				}
			}
			settings := bo.Settings{}
			if nodeXray.Settings != nil && *nodeXray.Settings != "" {
				if err = json.Unmarshal([]byte(*nodeXray.Settings), &settings); err != nil {
					logrus.Errorln(fmt.Sprintf("SystemVo JSON反转失败 err: %v", err))
					vo.Fail(constant.SysError, c)
					return
				}
			}
			if *nodeXray.Protocol == "vmess" {
				vmess := bo.Vmess{}
				vmess.Name = item.Name
				vmess.Server = item.Ip
				vmess.Port = item.Port
				vmess.VmessType = "vmess"
				vmess.Uuid = util.GenerateUUID(pass)
				vmess.AlterId = 0
				if settings.Encryption != "none" {
					vmess.Cipher = "auto"
				} else {
					vmess.Cipher = "none"
				}
				vmess.Udp = true
				vmess.Network = streamSettings.Network
				if streamSettings.Security == "tls" {
					vmess.Tls = true
				} else {
					vmess.Tls = false
				}
				if streamSettings.Network == "ws" {
					vmess.WsOpts.Path = streamSettings.WsSettings.Path
					vmess.WsOpts.WsOptsHeaders.Host = streamSettings.WsSettings.Host
				}
				ClashConfigInterface = append(ClashConfigInterface, vmess)
				proxies = append(proxies, item.Name)
			} else if *nodeXray.Protocol == "trojan" {
				trojan := bo.Trojan{}
				trojan.Name = item.Name
				trojan.Server = item.Ip
				trojan.Port = item.Port
				trojan.TrojanType = "trojan"
				trojan.Password = pass
				trojan.Udp = true
				trojan.Network = streamSettings.Network
				if streamSettings.Network == "ws" {
					trojan.WsOpts.Path = streamSettings.WsSettings.Path
					trojan.WsOpts.WsOptsHeaders.Host = streamSettings.WsSettings.Host
				}
				ClashConfigInterface = append(ClashConfigInterface, trojan)
				proxies = append(proxies, item.Name)
			}

		} else if item.NodeTypeId == 2 {
			nodeXrayTrojanGo, err := service.SelectNodeTrojanGoById(&item.NodeSubId)
			if err != nil {
				vo.Fail(err.Error(), c)
				return
			}
			trojanGo := bo.TrojanGo{}
			trojanGo.Name = item.Name
			trojanGo.Server = item.Ip
			trojanGo.Port = item.Port
			trojanGo.TrojanType = "trojan"
			trojanGo.Password = pass
			trojanGo.Udp = true
			trojanGo.SNI = *nodeXrayTrojanGo.Sni
			if *nodeXrayTrojanGo.WebsocketEnable == 1 {
				trojanGo.Network = "ws"
				trojanGo.WsOpts.Path = *nodeXrayTrojanGo.WebsocketPath
				trojanGo.WsOpts.WsOptsHeaders.Host = *nodeXrayTrojanGo.WebsocketHost
			}
			ClashConfigInterface = append(ClashConfigInterface, trojanGo)
			proxies = append(proxies, item.Name)
		}
	}
	proxyGroups := make([]bo.ProxyGroup, 0)
	proxyGroup := bo.ProxyGroup{
		Name:      "PROXY",
		ProxyType: "select",
		Proxies:   proxies,
	}
	proxyGroups = append(proxyGroups, proxyGroup)
	clashConfig.ProxyGroups = proxyGroups
	clashConfig.Proxies = ClashConfigInterface

	clashConfigYaml, err := yaml.Marshal(&clashConfig)
	if err != nil {
		vo.Fail(constant.SysError, c)
		return
	}

	result := fmt.Sprintf(`%s
%s`, string(clashConfigYaml), constant.ClashRules)

	c.Header("content-disposition", fmt.Sprintf("attachment; filename=%s.yaml", *account.Username))
	c.Header("profile-update-interval", "12")
	c.Header("subscription-userinfo", userInfo)
	c.String(200, result)
	return
}
