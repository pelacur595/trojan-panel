package service

import (
	"trojan-panel/dao"
	"trojan-panel/module/vo"
)

func SelectFileTaskPage(taskType *uint, pageNum *uint, pageSize *uint) (*vo.FileTaskPageVo, error) {
	return dao.SelectFileTaskPage(taskType, pageNum, pageSize)
}

func DeleteFileTaskById(id *uint) error {
	return dao.DeleteFileTaskById(id)
}
