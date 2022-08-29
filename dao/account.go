package dao

import (
	"errors"
	"fmt"
	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/sirupsen/logrus"
	"trojan/module"
	"trojan/module/constant"
	"trojan/module/vo"
	"trojan/util"
)

func SelectAccountById(id *uint) (*module.Account, error) {
	var account module.Account

	where := map[string]interface{}{"id": *id}
	selectFields := []string{"id", "username", "role_id", "email", "expire_time", "deleted", "quota",
		"download", "upload"}
	buildSelect, values, err := builder.BuildSelect("account", where, selectFields)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	if err = db.QueryRow(buildSelect, values...).Scan(&account); err == scanner.ErrEmptyResult {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.UnauthorizedError)
	} else if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	return &account, nil
}

func CreateAccount(account *module.Account) error {
	// 密码加密
	encryPass, err := util.AesEncode(*account.Pass)
	if err != nil {
		return err
	}

	accountCreate := map[string]interface{}{
		"username": *account.Username,
		"`pass`":   encryPass,
	}
	if account.RoleId != nil {
		accountCreate["`role_id`"] = *account.RoleId
	}
	if account.Email != nil && *account.Email != "" {
		accountCreate["`email`"] = *account.Email
	}
	if account.ExpireTime != nil {
		accountCreate["expire_time"] = *account.ExpireTime
	}
	if account.Deleted != nil {
		accountCreate["deleted"] = *account.Deleted
	}
	if account.Quota != nil {
		accountCreate["`quota`"] = *account.Quota
	}
	if account.IpLimit != nil {
		accountCreate["ip_limit"] = *account.IpLimit
	}
	if account.UploadSpeedLimit != nil {
		accountCreate["upload_speed_limit"] = *account.UploadSpeedLimit
	}
	if account.DownloadSpeedLimit != nil {
		accountCreate["download_speed_limit"] = *account.DownloadSpeedLimit
	}
	var data []map[string]interface{}
	data = append(data, accountCreate)

	buildInsert, values, err := builder.BuildInsert("account", data)
	if err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}
	if _, err = db.Exec(buildInsert, values...); err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}
	return nil
}

func CountAccountByUsername(username *string) (int, error) {
	var count int

	where := map[string]interface{}{
		"username": *username,
	}
	selectFields := []string{"count(1)"}
	buildSelect, values, err := builder.BuildSelect("account", where, selectFields)
	if err != nil {
		logrus.Errorln(err.Error())
		return 0, errors.New(constant.SysError)
	}

	if err = db.QueryRow(buildSelect, values...).Scan(&count); err != nil {
		logrus.Errorln(err.Error())
		return 0, errors.New(constant.SysError)
	}
	return count, nil
}

func SelectAccountPage(queryUsername *string, pageNum *uint, pageSize *uint) (*vo.AccountPageVo, error) {
	var (
		total    uint
		accounts []module.Account
	)

	// 查询总数
	var whereCount = map[string]interface{}{}
	if queryUsername != nil && *queryUsername != "" {
		whereCount["username like"] = fmt.Sprintf("%%%s%%", *queryUsername)
	}
	selectFieldsCount := []string{"count(1)"}
	buildSelect, values, err := builder.BuildSelect("account", whereCount, selectFieldsCount)
	if err = db.QueryRow(buildSelect, values...).Scan(&total); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	// 分页查询
	where := map[string]interface{}{
		"_orderby": "role_id,create_time desc",
		"_limit":   []uint{(*pageNum - 1) * *pageSize, *pageSize}}
	if queryUsername != nil && *queryUsername != "" {
		where["username like"] = fmt.Sprintf("%%%s%%", *queryUsername)
	}
	selectFields := []string{"id", "username", "role_id", "email", "expire_time", "deleted",
		"quota", "upload", "download", "create_time"}
	selectSQL, values, err := builder.BuildSelect("account", where, selectFields)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	rows, err := db.Query(selectSQL, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	defer rows.Close()

	if err = scanner.Scan(rows, &accounts); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	var accountVos = make([]vo.AccountVo, 0)
	for _, item := range accounts {
		accountVos = append(accountVos, vo.AccountVo{
			Id:         *item.Id,
			Username:   *item.Username,
			RoleId:     *item.RoleId,
			Email:      *item.Email,
			ExpireTime: *item.ExpireTime,
			Deleted:    *item.Deleted,
			Quota:      *item.Quota,
			Download:   *item.Download,
			Upload:     *item.Upload,
			CreateTime: *item.CreateTime,
		})
	}

	accountPageVo := vo.AccountPageVo{
		BaseVoPage: vo.BaseVoPage{
			PageNum:  *pageNum,
			PageSize: *pageSize,
			Total:    total,
		},
		AccountVos: accountVos,
	}
	return &accountPageVo, nil
}

func DeleteAccountById(id *uint) error {
	buildDelete, values, err := builder.BuildDelete("account", map[string]interface{}{"id": *id})
	if err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}

	if _, err = db.Exec(buildDelete, values...); err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}
	return nil
}

func SelectAccountByUsernameAndPass(username *string, pass *string) (*module.Account, error) {
	var account module.Account

	encryPass, err := util.AesEncode(*pass)
	if err != nil {
		return nil, err
	}
	where := map[string]interface{}{"username": *username, "pass": encryPass}
	selectFields := []string{"id", "username", "role_id", "deleted"}
	buildSelect, values, err := builder.BuildSelect("account", where, selectFields)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	rows, err := db.Query(buildSelect, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	defer rows.Close()

	if err = scanner.Scan(rows, &account); err == scanner.ErrEmptyResult {
		return nil, errors.New(constant.UsernameOrPassError)
	} else if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	return &account, nil
}

func UpdateAccountProfile(oldPass *string, newPass *string, username *string, email *string) error {
	_, err := SelectAccountByUsernameAndPass(username, oldPass)
	if err != nil {
		return errors.New(constant.OriPassError)
	}

	where := map[string]interface{}{"username": *username}
	update := map[string]interface{}{}
	if oldPass != nil && *oldPass != "" && newPass != nil && *newPass != "" {
		encryPass, err := util.AesEncode(*newPass)
		if err != nil {
			return err
		}
		update["pass"] = encryPass
	}
	if email != nil && *email != "" {
		update["email"] = *email
	}
	if len(update) > 0 {
		buildUpdate, values, err := builder.BuildUpdate("account", where, update)
		if err != nil {
			logrus.Errorln(err.Error())
			return errors.New(constant.SysError)
		}

		if _, err = db.Exec(buildUpdate, values...); err != nil {
			logrus.Errorln(err.Error())
			return errors.New(constant.SysError)
		}
	}
	return nil
}

func UpdateAccountById(account *module.Account) error {
	where := map[string]interface{}{"id": *account.Id}
	update := map[string]interface{}{}
	if account.Pass != nil && *account.Pass != "" {
		encryPass, err := util.AesEncode(*account.Pass)
		if err != nil {
			return err
		}
		update["pass"] = encryPass
	}
	if account.RoleId != nil {
		update["role_id"] = *account.RoleId
	}
	if account.Email != nil {
		update["email"] = *account.Email
	}
	if account.ExpireTime != nil {
		update["expire_time"] = *account.ExpireTime
	}
	if account.RoleId != nil {
		update["deleted"] = *account.Deleted
	}
	if account.Quota != nil {
		update["quota"] = *account.Quota
	}
	if account.IpLimit != nil {
		update["ip_limit"] = *account.IpLimit
	}
	if account.UploadSpeedLimit != nil {
		update["upload_speed_limit"] = *account.UploadSpeedLimit
	}
	if account.DownloadSpeedLimit != nil {
		update["download_speed_limit"] = *account.DownloadSpeedLimit
	}

	if len(update) > 0 {
		buildUpdate, values, err := builder.BuildUpdate("account", where, update)
		if err != nil {
			logrus.Errorln(err.Error())
			return errors.New(constant.SysError)
		}

		if _, err = db.Exec(buildUpdate, values...); err != nil {
			logrus.Errorln(err.Error())
			return errors.New(constant.SysError)
		}
	}
	return nil
}

func AccountQRCode(id *uint) (string, error) {
	var account module.Account

	where := map[string]interface{}{"id": *id}
	selectFields := []string{"username", "pass"}
	buildSelect, values, err := builder.BuildSelect("account", where, selectFields)
	if err != nil {
		logrus.Errorln(err.Error())
		return "", errors.New(constant.SysError)
	}

	if err = db.QueryRow(buildSelect, values...).Scan(&account); err == scanner.ErrEmptyResult {
		return "", errors.New(constant.UnauthorizedError)
	} else if err != nil {
		logrus.Errorln(err.Error())
		return "", errors.New(constant.SysError)
	}

	if account.Username == nil || *account.Username == "" || account.Pass == nil || *account.Pass == "" {
		return "", errors.New(constant.NodeQRCodeError)
	}
	decodePass, err := util.AesDecode(*account.Pass)
	if err != nil {
		logrus.Errorln(err.Error())
		return "", errors.New(constant.SysError)
	}
	return fmt.Sprintf("%s%s", *account.Username, decodePass), nil
}

func UpdateAccountQuotaOrDownloadOrUploadOrDeletedByUsernames(usernames []string, quota *int, download *uint, upload *uint, deleted *uint) error {
	where := map[string]interface{}{"username in": usernames}

	update := map[string]interface{}{}
	if quota != nil {
		update["quota"] = *quota
	}
	if download != nil {
		update["download"] = *download
	}
	if upload != nil {
		update["upload"] = *upload
	}
	if deleted != nil {
		update["deleted"] = *deleted
	}
	if len(update) > 0 {
		buildUpdate, values, err := builder.BuildUpdate("account", where, update)
		if err != nil {
			logrus.Errorln(err.Error())
			return errors.New(constant.SysError)
		}

		if _, err = db.Exec(buildUpdate, values...); err != nil {
			logrus.Errorln(err.Error())
			return errors.New(constant.SysError)
		}
	}
	return nil
}

// 查询禁用或者过期的用户名
func SelectAccountUsernameByDeletedOrExpireTime() ([]string, error) {
	buildSelect, values, err := builder.NamedQuery("select username from account where (deleted = {{deleted}} or expire_time <= {{expire_time}}) and quota != 0",
		map[string]interface{}{"deleted": 1, "expire_time": util.NowMilli()})
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	rows, err := db.Query(buildSelect, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	defer rows.Close()

	var usernames []string
	if err = scanner.Scan(rows, &usernames); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	return usernames, nil
}

// 用于发邮件
func SelectAccountsByExpireTime(expireTime uint) ([]module.Account, error) {
	buildSelect, values, err := builder.NamedQuery("select username,email from account where expire_time <= {{expire_time}} and quota != 0",
		map[string]interface{}{"expire_time": expireTime})
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	rows, err := db.Query(buildSelect, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	defer rows.Close()

	var accounts []module.Account
	if err = scanner.Scan(rows, &accounts); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	return accounts, nil
}

// TrafficRank 流量排行 前15名
func TrafficRank() ([]vo.AccountTrafficRankVo, error) {
	buildSelect, values, err := builder.NamedQuery(`select username, upload + download as traffic_used
from account
where quota != 0 and username not like '%admin%'
order by traffic_used desc limit 15`, nil)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	rows, err := db.Query(buildSelect, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	defer rows.Close()

	var accountTrafficRankVos []vo.AccountTrafficRankVo
	result, err := scanner.ScanMap(rows)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	for _, record := range result {
		accountTrafficRankVos = append(accountTrafficRankVos, vo.AccountTrafficRankVo{
			Username:    fmt.Sprintf("%s", record["username"]),
			TrafficUsed: fmt.Sprintf("%s", record["traffic_used"]),
		})
	}
	return accountTrafficRankVos, nil
}

// UpdateAccountQuota 重设流量 设置角色为普通用户流量为0
func UpdateAccountQuota() error {
	where := map[string]interface{}{"role_id": 3}
	update := map[string]interface{}{"quota": 0}
	buildUpdate, values, err := builder.BuildUpdate("account", where, update)
	if err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}

	if _, err = db.Exec(buildUpdate, values...); err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}
	return nil
}
