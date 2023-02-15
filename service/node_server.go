package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"sync"
	"trojan-panel/core"
	"trojan-panel/dao"
	"trojan-panel/module"
	"trojan-panel/module/constant"
	"trojan-panel/module/dto"
	"trojan-panel/module/vo"
	"trojan-panel/util"
)

func SelectNodeServerById(id *uint) (*module.NodeServer, error) {
	return dao.SelectNodeServerById(id)
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
					status, ok := nodeMap.Load(ip)
					if ok {
						splitNodeServerVos[indexI][j].Status = status.(int)
					} else {
						var status = 0
						success, err := core.Ping(token, ip, grpcPort)
						if err != nil {
							status = -1
						} else {
							if success {
								status = 1
							}
						}
						splitNodeServerVos[indexI][j].Status = status
						nodeMap.Store(ip, status)
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

func NodeServerState(token string, nodeServerId *uint) (*core.NodeServerGroupVo, error) {
	nodeServer, err := dao.SelectNodeServerById(nodeServerId)
	if err != nil {
		return nil, err
	}
	nodeServerGroupVo, err := core.NodeServerState(token, *nodeServer.Ip, *nodeServer.GrpcPort)
	if err != nil {
		return nil, err
	}
	return nodeServerGroupVo, nil
}
