package service

import (
	"os"
	"trojan-panel/dao"
	"trojan-panel/dao/redis"
	"trojan-panel/module"
	"trojan-panel/module/constant"
	"trojan-panel/module/vo"
)

func SelectFileTaskPage(taskType *uint, accountUsername *string, pageNum *uint, pageSize *uint) (*vo.FileTaskPageVo, error) {
	return dao.SelectFileTaskPage(taskType, accountUsername, pageNum, pageSize)
}

func DeleteFileTaskById(id *uint) error {
	mutex, err := redis.RsLock(constant.DeleteFileTaskByIdLock)
	if err != nil {
		return err
	}
	fileTask, err := dao.SelectFileTaskById(id)
	if err != nil {
		return err
	}
	if err = os.Remove(*fileTask.Path); err != nil {
		return err
	}
	if err := dao.DeleteFileTaskById(id); err != nil {
		return err
	}
	redis.RsUnLock(mutex)
	return nil
}

func SelectFileTaskById(id *uint) (*module.FileTask, error) {
	return dao.SelectFileTaskById(id)
}
