package service

import (
	"trojan-panel/dao"
	"trojan-panel/module"
	"trojan-panel/module/vo"
)

func SelectEmailRecordPage(queryToEmail *string, queryState *int, pageNum *uint, pageSize *uint) (*vo.EmailRecordPageVo, error) {
	return dao.SelectEmailRecordPage(queryToEmail, queryState, pageNum, pageSize)
}

func CreateEmailRecord(emailRecord module.EmailRecord) (uint, error) {
	return dao.CreateEmailRecord(emailRecord)
}
