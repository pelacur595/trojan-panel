package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"sync"
	"time"
	"trojan-panel/core"
	"trojan-panel/dao"
	"trojan-panel/module"
	"trojan-panel/module/constant"
	"trojan-panel/module/dto"
	"trojan-panel/module/vo"
	"trojan-panel/util"
)

func SelectNodeServerById(id *uint) (*module.NodeServer, error) {
	return dao.SelectNodeServer(map[string]interface{}{"id": *id})
}

func CreateNodeServer(nodeServer *module.NodeServer) error {
	count, err := dao.CountNodeServerByName(nil, nodeServer.Name)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(constant.NodeServerNameExist)
	}
	return dao.CreateNodeServer(nodeServer)
}

func SelectNodeServerPage(queryName *string, queryIp *string, pageNum *uint, pageSize *uint, c *gin.Context) (*vo.NodeServerPageVo, error) {
	nodeServerPage, total, err := dao.SelectNodeServerPage(queryName, queryIp, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	nodeServerVos := make([]vo.NodeServerVo, 0)
	for _, item := range *nodeServerPage {
		nodeServerVo := vo.NodeServerVo{
			Id:         *item.Id,
			Name:       *item.Name,
			Ip:         *item.Ip,
			GrpcPort:   *item.GrpcPort,
			CreateTime: *item.CreateTime,
		}
		nodeServerVos = append(nodeServerVos, nodeServerVo)
	}

	account := util.GetCurrentAccount(c)
	if util.IsAdmin(account.Roles) {
		token := util.GetToken(c)
		splitNodeServerVos := util.SplitArr(nodeServerVos, 2)
		var nodeMap sync.Map
		var wg sync.WaitGroup
		for i := range splitNodeServerVos {
			indexI := i
			wg.Add(1)
			go func() {
				for j := range splitNodeServerVos[indexI] {
					var ip = splitNodeServerVos[indexI][j].Ip
					var grpcPort = splitNodeServerVos[indexI][j].GrpcPort
					nodeMapValue, ok := nodeMap.Load(ip)
					if ok {
						nodeServerVo := nodeMapValue.(vo.NodeServerVo)
						splitNodeServerVos[indexI][j].Status = nodeServerVo.Status
						splitNodeServerVos[indexI][j].TrojanPanelCoreVersion = nodeServerVo.TrojanPanelCoreVersion
					} else {
						var nodeServerState int
						var trojanPanelCoreVersion string
						stateVo, err := core.GetNodeServerState(token, ip, grpcPort)
						if err != nil {
							nodeServerState = 0
						} else {
							nodeServerState = 1
							trojanPanelCoreVersion = stateVo.GetVersion()
						}
						splitNodeServerVos[indexI][j].Status = nodeServerState
						splitNodeServerVos[indexI][j].TrojanPanelCoreVersion = trojanPanelCoreVersion
						nodeMap.Store(ip, splitNodeServerVos[indexI][j])
					}
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}

	nodeServerPageVo := vo.NodeServerPageVo{
		BaseVoPage: vo.BaseVoPage{
			PageNum:  *pageNum,
			PageSize: *pageSize,
			Total:    total,
		},
		NodeServers: nodeServerVos,
	}
	return &nodeServerPageVo, nil
}

func DeleteNodeServerById(id *uint) error {
	count, err := dao.CountNodeByNameAndNodeServerId(nil, nil, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(constant.NodeServerDeletedError)
	}

	return dao.DeleteNodeServerById(id)
}

func UpdateNodeServerById(dto *dto.NodeServerUpdateDto) error {
	count, err := dao.CountNodeByNameAndNodeServerId(nil, nil, dto.Id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(constant.NodeServerDeletedError)
	}

	count, err = dao.CountNodeServerByName(dto.Id, dto.Name)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(constant.NodeServerNameExist)
	}

	nodeServer := module.NodeServer{
		Id:       dto.Id,
		Ip:       dto.Ip,
		Name:     dto.Name,
		GrpcPort: dto.GrpcPort,
	}
	return dao.UpdateNodeServerById(&nodeServer)
}

func CountNodeServer() (int, error) {
	return dao.CountNodeServer()
}

func CountNodeServerByName(id *uint, queryName *string) (int, error) {
	return dao.CountNodeServerByName(id, queryName)
}

func SelectNodeServerList(dto *dto.NodeServerDto) ([]vo.NodeServerListVo, error) {
	nodeServerList, err := dao.SelectNodeServerList(dto.Ip, dto.Name)
	if err != nil {
		return nil, err
	}
	nodeServerListVos := make([]vo.NodeServerListVo, 0)
	for _, item := range nodeServerList {
		nodeServerVo := vo.NodeServerListVo{
			Id:   *item.Id,
			Name: *item.Name,
		}
		nodeServerListVos = append(nodeServerListVos, nodeServerVo)

	}
	return nodeServerListVos, nil
}

func GetNodeServerInfo(token string, nodeServerId *uint) (*core.NodeServerInfoVo, error) {
	nodeServer, err := dao.SelectNodeServer(map[string]interface{}{"id": *nodeServerId})
	if err != nil {
		return nil, err
	}
	nodeServerInfoVo, err := core.GetNodeServerInfo(token, *nodeServer.Ip, *nodeServer.GrpcPort)
	if err != nil {
		return nil, err
	}
	return nodeServerInfoVo, nil
}

func ExportNodeServer(accountId uint, accountUsername string) error {
	fileName := fmt.Sprintf("nodeServerExport-%s.csv", time.Now().Format("20060102150405"))
	filePath := fmt.Sprintf("%s/%s", constant.ExcelPath, fileName)

	var fileTaskType uint = constant.TaskTypeNodeServerExport
	var fileTaskStatus = constant.TaskDoing
	fileTask := module.FileTask{
		Name:            &fileName,
		Path:            &filePath,
		Type:            &fileTaskType,
		Status:          &fileTaskStatus,
		AccountId:       &accountId,
		AccountUsername: &accountUsername,
	}
	fileTaskId, err := dao.CreateFileTask(&fileTask)
	if err != nil {
		return err
	}

	go func() {
		var mutex sync.Mutex
		defer mutex.Unlock()
		if mutex.TryLock() {
			var fail = constant.TaskFail
			var success = constant.TaskSuccess
			fileTask := module.FileTask{
				Id:     &fileTaskId,
				Status: &fail,
			}

			var data [][]string
			titles := []string{"ip", "name", "grpc_port", "create_time"}
			data = append(data, titles)
			// 查询所有需要导出数据
			nodeServerExportVo, err := dao.SelectNodeServerAll()
			if err != nil {
				logrus.Errorf("ExportNodeServer SelectNodeServerAll err: %v", err)
			}
			for _, item := range nodeServerExportVo {
				element := []string{item.Ip, item.Name, item.GrpcPort, item.CreateTime}
				data = append(data, element)
			}
			if err = util.ExportCsv(filePath, data); err != nil {
				logrus.Errorf("ExportNodeServer ExportCsv err: %v", err)
			} else {
				fileTask.Status = &success
			}

			// 更新文件任务状态
			if err = dao.UpdateFileTaskById(&fileTask); err != nil {
				logrus.Errorf("ExportNodeServer UpdateFileTaskById err: %v", err)
			}
		}
	}()
	return nil
}

func ImportNodeServer(cover uint, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	reader := csv.NewReader(src)

	// 读取表头
	titlesRead, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return errors.New(constant.CsvRowNotEnough)
		}
		logrus.Errorf("ImportNodeServer read csv titles err: %s", err.Error())
	}
	titles := []string{"ip", "name", "grpc_port"}
	// 必须以titles作为表头
	if !util.ArraysEqualPrefix(titles, titlesRead) {
		return errors.New(constant.CsvTitleError)
	}
	// data 变量中存储CSV文件中的数据
	var data [][]string
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			logrus.Errorf("ImportNodeServer read csv record err: %s", err.Error())
		}
		data = append(data, record)
	}
	// 在这里可以处理数据并将其存储到数据库中 todo 这里可能存在性能问题
	for _, item := range data {
		if err = dao.CreateOrUpdateNodeServer(item, cover); err != nil {
			return err
		}
	}
	return nil
}
