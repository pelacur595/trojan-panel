package service

import (
	"os"
	"sync"
	"trojan-panel/dao"
	"trojan-panel/module"
	"trojan-panel/module/vo"
)

func SelectFileTaskPage(taskType *uint, pageNum *uint, pageSize *uint) (*vo.FileTaskPageVo, error) {
	return dao.SelectFileTaskPage(taskType, pageNum, pageSize)
}

func DeleteFileTaskById(id *uint) error {
	var mutex sync.Mutex
	defer mutex.TryLock()
	if mutex.TryLock() {
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
	}
	return nil
}

func SelectFileTaskById(id *uint) (*module.FileTask, error) {
	return dao.SelectFileTaskById(id)
}
