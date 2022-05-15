package dao

import (
	"crypto/sha256"
	"encoding/base64"
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

func SelectUserById(id *uint) (*vo.UsersVo, error) {
	var user module.Users

	where := map[string]interface{}{"id": *id}
	selectFields := []string{"id", "role_id", "username", "email", "quota", "FLOOR(upload/1024/1024) upload", "download,deleted,expire_time"}
	buildSelect, values, err := builder.BuildSelect("users", where, selectFields)
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

	err = scanner.Scan(rows, &user)
	if err == scanner.ErrEmptyResult {
		return nil, errors.New(constant.UnauthorizedError)
	} else if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	usersVo := vo.UsersVo{
		Id:         *user.Id,
		Quota:      *user.Quota,
		Username:   *user.Username,
		Email:      *user.Email,
		RoleId:     *user.RoleId,
		Upload:     *user.Upload,
		Download:   *user.Download,
		Deleted:    *user.Deleted,
		ExpireTime: *user.ExpireTime,
	}

	return &usersVo, nil
}

func CreateUser(users *module.Users) error {
	// 密码加密
	encryPass := base64.StdEncoding.EncodeToString([]byte(*users.Pass))
	encryPassword := fmt.Sprintf("%x", sha256.Sum224([]byte(fmt.Sprintf("%s&%s", *users.Username, *users.Pass))))

	var data []map[string]interface{}
	user := map[string]interface{}{
		"`password`": encryPassword,
		"username":   *users.Username,
		"`pass`":     encryPass,
	}
	if users.Quota != nil {
		user["`quota`"] = *users.Quota
	}
	if users.Email != nil && *users.Email != "" {
		user["`email`"] = *users.Email
	}
	if users.RoleId != nil {
		user["`role_id`"] = *users.RoleId
	}
	if users.Deleted != nil {
		user["deleted"] = *users.Deleted
	}
	if users.ExpireTime != nil {
		user["expire_time"] = *users.ExpireTime
	}
	if users.IpLimit != nil {
		user["ip_limit"] = *users.IpLimit
	}
	if users.UploadSpeedLimit != nil {
		user["upload_speed_limit"] = *users.UploadSpeedLimit
	}
	if users.DownloadSpeedLimit != nil {
		user["download_speed_limit"] = *users.DownloadSpeedLimit
	}
	data = append(data, user)

	buildInsert, values, err := builder.BuildInsert("users", data)
	if err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}
	if _, err := db.Exec(buildInsert, values...); err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}
	return nil
}

func CountUserByUsername(username *string) (int, error) {
	var count int

	where := map[string]interface{}{}
	if username != nil {
		where["username"] = *username
	}
	selectFields := []string{"count(1)"}
	buildSelect, values, err := builder.BuildSelect("users", where, selectFields)
	if err != nil {
		logrus.Errorln(err.Error())
		return 0, errors.New(constant.SysError)
	}

	if err := db.QueryRow(buildSelect, values...).Scan(&count); err != nil {
		logrus.Errorln(err.Error())
		return 0, errors.New(constant.SysError)
	}
	return count, nil
}

func SelectUserPage(queryUsername *string, pageNum *uint, pageSize *uint) (*vo.UsersPageVo, error) {
	var (
		total uint
		users []module.Users
	)

	// 查询总数
	var whereCount = map[string]interface{}{}
	if queryUsername != nil && *queryUsername != "" {
		whereCount["username like"] = fmt.Sprintf("%%%s%%", *queryUsername)
	}
	selectFieldsCount := []string{"count(1)"}
	buildSelect, values, err := builder.BuildSelect("users", whereCount, selectFieldsCount)
	if err := db.QueryRow(buildSelect, values...).Scan(&total); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	// 分页查询
	offset := (*pageNum - 1) * *pageSize
	where := map[string]interface{}{
		"_orderby": "role_id,create_time desc",
		"_limit":   []uint{offset, *pageSize}}
	if queryUsername != nil && *queryUsername != "" {
		where["username like"] = fmt.Sprintf("%%%s%%", *queryUsername)
	}
	selectFields := []string{"id", "role_id", "username", "quota",
		"upload", "download", "deleted", "expire_time", "email", "create_time"}
	selectSQL, values, err := builder.BuildSelect("users", where, selectFields)
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

	if err := scanner.Scan(rows, &users); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	var userVos []vo.UsersVo
	for _, item := range users {
		userVos = append(userVos, vo.UsersVo{
			Id:         *item.Id,
			Quota:      util.ToMB(*item.Quota),
			Download:   *item.Download,
			Upload:     *item.Upload,
			Username:   *item.Username,
			RoleId:     *item.RoleId,
			CreateTime: *item.CreateTime,
			Deleted:    *item.Deleted,
			ExpireTime: *item.ExpireTime,
			Email:      *item.Email,
		})
	}

	usersPageVo := vo.UsersPageVo{
		BaseVoPage: vo.BaseVoPage{
			PageNum:  *pageNum,
			PageSize: *pageSize,
			Total:    total,
		},
		Users: userVos,
	}
	return &usersPageVo, nil
}

func DeleteUserById(id *uint) error {
	buildDelete, values, err := builder.BuildDelete("users", map[string]interface{}{"id": *id})
	if err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}

	if _, err := db.Exec(buildDelete, values...); err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}
	return nil
}

func SelectUserByUsernameAndPass(username *string, pass *string) (*vo.UsersVo, error) {
	var user module.Users

	encryPass := base64.StdEncoding.EncodeToString([]byte(*pass))
	where := map[string]interface{}{"username": *username, "pass": encryPass}
	selectFileds := []string{"id", "role_id", "username", "deleted"}
	buildSelect, values, err := builder.BuildSelect("users", where, selectFileds)
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

	err = scanner.Scan(rows, &user)
	if err == scanner.ErrEmptyResult {
		return nil, errors.New(constant.UsernameOrPassError)
	} else if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	usersVo := vo.UsersVo{
		Id:       *user.Id,
		Username: *user.Username,
		RoleId:   *user.RoleId,
		Deleted:  *user.Deleted,
	}
	return &usersVo, nil
}

func UpdateUserProfile(oldPass *string, newPass *string, username *string, email *string) error {
	_, err := SelectUserByUsernameAndPass(username, oldPass)
	if err != nil {
		return errors.New(constant.OriPassError)
	}
	encryPass := base64.StdEncoding.EncodeToString([]byte(*newPass))
	encryPassword := sha256.Sum224([]byte(fmt.Sprintf("%s&%s", *username, *newPass)))

	where := map[string]interface{}{"username": *username}
	update := map[string]interface{}{
		"`pass`":   encryPass,
		"password": fmt.Sprintf("%x", encryPassword),
	}
	if email != nil && *email != "" {
		update["email"] = email
	}
	buildUpdate, values, err := builder.BuildUpdate("users", where, update)
	if err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}

	if _, err := db.Exec(buildUpdate, values...); err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}
	return nil
}

func UpdateUserById(users *module.Users) error {
	where := map[string]interface{}{"id": *users.Id}
	update := map[string]interface{}{}
	if users.Pass != nil && *users.Pass != "" {
		encryPass := base64.StdEncoding.EncodeToString([]byte(*users.Pass))
		encryPassword := sha256.Sum224([]byte(fmt.Sprintf("%s&%s", *users.Username, *users.Pass)))
		update["`pass`"] = encryPass
		update["`password`"] = fmt.Sprintf("%x", encryPassword)
	}
	if users.Quota != nil {
		update["quota"] = *users.Quota
	}
	if users.Email != nil {
		update["email"] = *users.Email
	}
	if users.RoleId != nil {
		update["role_id"] = *users.RoleId
	}
	if users.RoleId != nil {
		update["deleted"] = *users.Deleted
	}
	if users.ExpireTime != nil {
		update["expire_time"] = *users.ExpireTime
	}
	if users.IpLimit != nil {
		update["ip_limit"] = *users.IpLimit
	}
	if users.UploadSpeedLimit != nil {
		update["upload_speed_limit"] = *users.UploadSpeedLimit
	}
	if users.DownloadSpeedLimit != nil {
		update["download_speed_limit"] = *users.DownloadSpeedLimit
	}

	if len(update) > 0 {
		buildUpdate, values, err := builder.BuildUpdate("users", where, update)
		if err != nil {
			logrus.Errorln(err.Error())
			return errors.New(constant.SysError)
		}

		_, err = db.Exec(buildUpdate, values...)
		if err != nil {
			logrus.Errorln(err.Error())
			return errors.New(constant.SysError)
		}
	}
	return nil
}

func UserQRCode(id *uint) (string, error) {
	var user module.Users
	where := map[string]interface{}{"id": *id}
	selectFields := []string{"username", "pass"}
	buildSelect, values, err := builder.BuildSelect("users", where, selectFields)
	if err != nil {
		logrus.Errorln(err.Error())
		return "", errors.New(constant.SysError)
	}

	rows, err := db.Query(buildSelect, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return "", errors.New(constant.SysError)
	}
	defer rows.Close()

	err = scanner.Scan(rows, &user)
	if err == scanner.ErrEmptyResult {
		return "", errors.New(constant.UnauthorizedError)
	} else if err != nil {
		logrus.Errorln(err.Error())
		return "", errors.New(constant.SysError)
	}
	if user.Username == nil || *user.Username == "" || user.Pass == nil || *user.Pass == "" {
		return "", errors.New(constant.NodeQRCodeError)
	}
	decodePass, err := base64.StdEncoding.DecodeString(*user.Pass)
	if err != nil {
		logrus.Errorln(err.Error())
		return "", errors.New(constant.SysError)
	}
	return fmt.Sprintf("%s&%s", *user.Username, string(decodePass)), nil
}

func UpdateUserQuotaOrDownloadOrUploadOrDeletedByUsernames(usernames []string, quota *int, download *uint, upload *uint, deleted *uint) error {
	where := map[string]interface{}{"username in": usernames}

	update := map[string]interface{}{}
	if quota != nil {
		update["quota"] = quota
	}
	if download != nil {
		update["download"] = download
	}
	if upload != nil {
		update["upload"] = upload
	}
	if deleted != nil {
		update["deleted"] = deleted
	}
	if len(update) > 0 {
		buildUpdate, values, err := builder.BuildUpdate("users", where, update)
		if err != nil {
			logrus.Errorln(err.Error())
			return errors.New(constant.SysError)
		}

		_, err = db.Exec(buildUpdate, values...)
		if err != nil {
			logrus.Errorln(err.Error())
			return errors.New(constant.SysError)
		}
	}
	return nil
}

// 查询禁用或者过期的用户名
func SelectUsernameByDeletedOrExpireTime() ([]string, error) {
	buildSelect, values, err := builder.NamedQuery("select username from users where (deleted = {{deleted}} or expire_time <= {{expire_time}}) and quota != 0",
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
	if err := scanner.Scan(rows, &usernames); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	return usernames, nil
}

func SelectUsersEmailByExpireTime(day uint) ([]module.Users, error) {
	buildSelect, values, err := builder.NamedQuery("select username,email from users where expire_time <= {{expire_time}} and quota != 0",
		map[string]interface{}{"expire_time": day})
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

	var users []module.Users
	if err := scanner.Scan(rows, &users); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	return users, nil
}

func SelectUserPasswordByUsernameOrId(id *uint, username *string) (string, error) {
	where := map[string]interface{}{}
	if id != nil {
		where["id"] = *id
	}
	if username != nil && *username != "" {
		where["username"] = *username
	}
	selectFileds := []string{"username"}
	buildSelect, values, err := builder.BuildSelect("`users`", where, selectFileds)
	if err != nil {
		logrus.Errorln(err.Error())
		return "", errors.New(constant.SysError)
	}
	var password string
	if err := db.QueryRow(buildSelect, values...).Scan(&password); err != nil {
		return "", errors.New(constant.SysError)
	}
	return password, err
}

func TrafficRank() ([]vo.UsersTrafficRankVo, error) {
	buildSelect, values, err := builder.NamedQuery(`select username, upload + download as traffic_used
from users
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

	var usersTrafficRankVos []vo.UsersTrafficRankVo
	result, err := scanner.ScanMap(rows)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	for _, record := range result {
		usersTrafficRankVos = append(usersTrafficRankVos, vo.UsersTrafficRankVo{
			Username:    fmt.Sprintf("%s", record["username"]),
			TrafficUsed: fmt.Sprintf("%s", record["traffic_used"]),
		})
	}
	return usersTrafficRankVos, nil
}
