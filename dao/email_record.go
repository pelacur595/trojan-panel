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
)

func SelectEmailRecordPage(queryToEmail *string, pageNum *uint, pageSize *uint) (*vo.EmailRecordPageVo, error) {
	var (
		total        uint
		emailRecords []module.EmailRecord
	)

	// 查询总数
	var whereCount = map[string]interface{}{}
	if queryToEmail != nil && *queryToEmail != "" {
		whereCount["to_email like"] = fmt.Sprintf("%%%s%%", *queryToEmail)
	}

	selectFieldsCount := []string{"count(1)"}
	buildSelect, values, err := builder.BuildSelect("email_record", whereCount, selectFieldsCount)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	if err := db.QueryRow(buildSelect, values...).Scan(&total); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	// 分页查询
	offset := (*pageNum - 1) * *pageSize
	where := map[string]interface{}{
		"_orderby": "create_time desc",
		"_limit":   []uint{offset, *pageSize}}
	if queryToEmail != nil && *queryToEmail != "" {
		where["to_email like"] = fmt.Sprintf("%%%s%%", *queryToEmail)
	}
	selectFields := []string{"id", "`to_email`", "subject", "content",
		"state", "create_time"}
	selectSQL, values, err := builder.BuildSelect("email_record", where, selectFields)
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

	if err = scanner.Scan(rows, &emailRecords); err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	var emailRecordVos []vo.EmailRecordVo
	for _, item := range emailRecords {
		emailRecordVos = append(emailRecordVos, vo.EmailRecordVo{
			Id:         *item.Id,
			ToEmail:    *item.ToEmail,
			Subject:    *item.Subject,
			Content:    *item.Content,
			State:      *item.State,
			CreateTime: *item.CreateTime,
		})
	}

	emailRecordPageVo := vo.EmailRecordPageVo{
		BaseVoPage: vo.BaseVoPage{
			PageNum:  *pageNum,
			PageSize: *pageSize,
			Total:    total,
		},
		EmailRecordVos: emailRecordVos,
	}
	return &emailRecordPageVo, nil
}

func CreateEmailRecord(emailRecord *module.EmailRecord) error {
	var data []map[string]interface{}
	data = append(data, map[string]interface{}{
		"to_email": *emailRecord.ToEmail,
		"subject":  *emailRecord.Subject,
		"content":  *emailRecord.Content,
		"state":    *emailRecord.State,
	})

	buildInsert, values, err := builder.BuildInsert("email_record", data)
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
