package service

import (
	"trojan/dao"
	"trojan/module"
	"trojan/module/vo"
)

func SelectEmailRecordPage(queryToEmail *string, queryState *int, pageNum *uint, pageSize *uint) (*vo.EmailRecordPageVo, error) {
	emailRecordPageVo, err := dao.SelectEmailRecordPage(queryToEmail, queryState, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	return emailRecordPageVo, nil
}

func CreateEmailRecord(emailRecord module.EmailRecord) (uint, error) {
	id, err := dao.CreateEmailRecord(emailRecord)
	if err != nil {
		return 0, err
	}
	return id, nil
}
