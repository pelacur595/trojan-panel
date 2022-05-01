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
	selectFields := []string{"id", "role_id", "username", "email", "FLOOR(quota/1024/1024) quota", "FLOOR(upload/1024/1024) upload", "FLOOR(download/1024/1024) download,deleted,expire_time"}
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
	data = append(data, map[string]interface{}{
		"`password`":  encryPassword,
		"`quota`":     *users.Quota,
		"username":    *users.Username,
		"`pass`":      encryPass,
		"email":       users.Email,
		"`role_id`":   *users.RoleId,
		"deleted":     *users.Deleted,
		"expire_time": *users.ExpireTime,
	})

	buildInsert, values, err := builder.BuildInsert("users", data)
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

func CountUserByUsername(username *string) (int, error) {
	var count int
	where := map[string]interface{}{}
	if username != nil {
		where["username"] = *username
	}
	selectFields := []string{"count(1) count"}
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
	selectFields := []string{"id", "role_id", "username", "FLOOR(quota/1024/1024) quota",
		"FLOOR(upload/1024/1024) upload", "FLOOR(download/1024/1024) download", "deleted", "expire_time", "create_time"}
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

	if err = scanner.Scan(rows, &users); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	var userVos []vo.UsersVo
	for _, item := range users {
		userVos = append(userVos, vo.UsersVo{
			Id:         *item.Id,
			Quota:      *item.Quota,
			Download:   *item.Download,
			Upload:     *item.Upload,
			Username:   *item.Username,
			RoleId:     *item.RoleId,
			CreateTime: *item.CreateTime,
			Deleted:    *item.Deleted,
			ExpireTime: *item.ExpireTime,
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

func UpdateUserPassByUsername(oldPass *string, newPass *string, username *string) error {
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
	if users.Email != nil {
		update["email"] = *users.Email
	}
	if users.Quota != nil {
		update["quota"] = *users.Quota
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

func UpdateUserPasswordOrDeletedByUsernames(usernames []string, password *string, deleted *uint) error {
	where := map[string]interface{}{"username in": usernames}

	update := map[string]interface{}{}
	if password != nil {
		update["password"] = password
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

// 通过用户名查询连接密码
func SelectEncryPasswordByUsername(username *string) (string, error) {
	var pass *string
	where := map[string]interface{}{"username": *username}
	selectFields := []string{"`pass`"}
	buildSelect, values, err := builder.BuildSelect("users", where, selectFields)
	if err != nil {
		logrus.Errorln(err.Error())
		return "", errors.New(constant.SysError)
	}

	if err = db.QueryRow(buildSelect, values...).Scan(&pass); err != nil {
		logrus.Errorln(err.Error())
		return "", errors.New(constant.SysError)
	}
	if pass == nil || *pass == "" {
		return "", errors.New(constant.UnauthorizedError)
	}

	return fmt.Sprintf("%x", sha256.Sum224([]byte(fmt.Sprintf("%s&%s", *username, *pass)))), nil
}

// 查询禁用或者过期的用户名
func SelectUsernameByDeletedOrExpireTime() ([]string, error) {
	var usernames []string
	buildSelect, values, err := builder.NamedQuery("select username from users where (deleted = {{deleted}} or expire_time <= {{expire_time}}) and password != ''",
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

	scanMap, err := scanner.ScanMap(rows)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	for _, record := range scanMap {
		usernames = append(usernames, fmt.Sprintf("%s", record["username"]))
	}

	return usernames, nil
}
