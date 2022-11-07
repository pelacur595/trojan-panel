package service

import (
	"trojan-panel/dao"
	"trojan-panel/module"
)

func SelectNodeXrayById(id *uint) (*module.NodeXray, error) {
	return dao.SelectNodeXrayById(id)
}
