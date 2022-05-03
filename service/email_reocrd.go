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

func CreateEmailRecord(emailRecords []module.EmailRecord) error {
	if err := dao.CreateEmailRecord(emailRecords); err != nil {
		return err
	}
	return nil
}

func UpdateEmailRecordById(emailRecord *module.EmailRecord) error {
	if err := dao.UpdateEmailRecordById(emailRecord); err != nil {
		return err
	}
	return nil
}
