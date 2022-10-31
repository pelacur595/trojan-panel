package service

import (
	"trojan-panel/dao"
	"trojan-panel/module"
)

func SelectNodeTrojanGoById(id *uint) (*module.NodeTrojanGo, error) {
	return dao.SelectNodeTrojanGoById(id)
}
