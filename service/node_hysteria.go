package service

import (
	"trojan-panel/dao"
	"trojan-panel/module"
)

func SelectNodeHysteriaById(id *uint) (*module.NodeHysteria, error) {
	return dao.SelectNodeHysteriaById(id)
}
