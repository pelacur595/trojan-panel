package service

import (
	"fmt"
	"trojan/dao"
	"trojan/dao/redis"
	"trojan/module/vo"
)

func DeleteBlackListByIp(ip *string) error {
	if err := dao.DeleteBlackListByIp(ip); err != nil {
		return err
	}
	redis.Client.Key.Del("trojan-panel:black-list:%s", *ip)
	return nil
}

func CreateBlackList(ips []string) error {
	if err := dao.CreateBlackList(ips); err != nil {
		return err
	}
	kv := map[string]interface{}{}
	for _, ip := range ips {
		kv[fmt.Sprintf("trojan-panel:black-list:%s", ip)] = "in-black-list"
	}
	redis.Client.String.MSet(kv)
	return nil
}

func SelectBlackListPage(ip *string, pageNum *uint, pageSize *uint) (*vo.BlackListPageVo, error) {
	blackListPageVo, err := dao.SelectBlackListPage(ip, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	return blackListPageVo, nil
}
