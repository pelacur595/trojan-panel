package dao

import (
	"errors"
	"fmt"
	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/sirupsen/logrus"
	"trojan/module"
	"trojan/module/constant"
	"trojan/module/dto"
	"trojan/module/vo"
)

func SelectMenuListByRoleId(roleId *uint) ([]vo.MenuListVo, error) {
	var menuLists []module.MenuList
	buildSelect, values, err := builder.NamedQuery(
		"select ml.id, ml.name, ml.icon, ml.route, ml.parent_id from role_menu_list rml left join menu_list ml on rml.menu_list_id = ml.id where rml.role_id = {{role_id}}",
		map[string]interface{}{"role_id": *roleId})
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

	if err = scanner.Scan(rows, &menuLists); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	var menuListVos []vo.MenuListVo
	for _, item := range menuLists {
		menuListVos = append(menuListVos, vo.MenuListVo{
			Id:       *item.Id,
			Name:     *item.Name,
			Icon:     *item.Icon,
			Route:    *item.Route,
			ParentId: *item.ParentId,
		})
	}
	return menuListVos, nil
}

func SelectRoleList(roleDto dto.RoleDto) (*[]vo.RoleListVo, error) {
	var roles []module.Role

	where := map[string]interface{}{}
	where["_orderby"] = "create_time desc"
	if roleDto.Name != nil && *roleDto.Name != "" {
		where["name"] = *roleDto.Name
	}
	if roleDto.Desc != nil && *roleDto.Desc != "" {
		where["desc"] = *roleDto.Desc
	}
	selectFields := []string{"id", "`name`", "`desc`"}
	buildSelect, values, err := builder.BuildSelect("role", where, selectFields)
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

	if err = scanner.Scan(rows, &roles); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	var roleListVos []vo.RoleListVo
	for _, item := range roles {
		roleListVos = append(roleListVos, vo.RoleListVo{
			Id:   *item.Id,
			Name: *item.Name,
			Desc: *item.Desc,
		})
	}
	return &roleListVos, nil
}

func SelectRoleNameByParentId(id *uint, includeSelf bool) ([]string, error) {
	var roleNames []string
	roleVo, err := SelectRoleById(id)
	if err != nil {
		return nil, err
	}
	if includeSelf {
		if roleVo.Name != "" {
			roleNames = append(roleNames, roleVo.Name)
		}
	}
	buildSelect, values, err := builder.NamedQuery("select `name` from `role` where `path` like concat({{path}},'-','%')",
		map[string]interface{}{"path": fmt.Sprintf("%s%d", roleVo.Path, *id)})
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

	result, err := scanner.ScanMap(rows)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	for _, record := range result {
		roleNames = append(roleNames, fmt.Sprintf("%s", record["name"]))
	}
	return roleNames, nil
}

func SelectRoleById(id *uint) (*vo.RoleVo, error) {
	var role module.Role
	buildSelect, values, err := builder.NamedQuery("select id,`name`,`desc`,`path` from `role` where id = {{id}}",
		map[string]interface{}{"id": *id})
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

	err = scanner.Scan(rows, &role)
	if err == scanner.ErrEmptyResult {
		return nil, errors.New(constant.RoleNotExist)
	} else if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	roleVo := vo.RoleVo{
		Id:   *role.Id,
		Name: *role.Name,
		Desc: *role.Desc,
		Path: *role.Path,
	}
	return &roleVo, nil
}
