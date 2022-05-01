package core

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/p4gefau1t/trojan-go/api/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"os"
	"trojan/module/constant"
)

type trojanGoApi struct {
	ctx context.Context
}

func TrojanGoApi() *trojanGoApi {
	return &trojanGoApi{
		ctx: context.Background(),
	}
}

func apiClient(addr string) (service.TrojanServerServiceClient, *grpc.ClientConn, error) {
	pool := x509.NewCertPool()
	certBytes, err := os.ReadFile(fmt.Sprintf("%s%s.crt", constant.TrojanPanelTrojanGoCertDate, addr))
	if err != nil {
		return nil, nil, err
	}
	pool.AppendCertsFromPEM(certBytes)

	certificate, err := tls.LoadX509KeyPair(constant.TrojanPanelCrtFile, constant.TrojanPanelKeyFile)
	if err != nil {
		logrus.Errorf("%s 加载本机密钥和证书失败 err: %v\n", addr, err)
		return nil, nil, errors.New(constant.LoadKeyPairError)
	}
	creds := credentials.NewTLS(&tls.Config{
		ServerName:   addr,
		RootCAs:      pool,
		Certificates: []tls.Certificate{certificate},
	})
	conn, err := grpc.Dial(fmt.Sprintf("%s:10010", addr), grpc.WithTransportCredentials(creds))
	if err != nil {
		logrus.Errorf("%s grpc连接化失败 err: %v\n", addr, err)
		return nil, nil, errors.New(constant.GrpcError)
	}
	return service.NewTrojanServerServiceClient(conn), conn, nil
}

func (t *trojanGoApi) OnLine(ip string) (int, error) {
	users, err := t.ListUsers(ip)
	if err != nil {
		return 0, err
	}
	var num = 0
	for _, user := range users {
		if user.SpeedCurrent != nil {
			num++
		}
	}
	return num, nil
}

// 查询节点上的所有用户
func (t *trojanGoApi) ListUsers(ip string) ([]*service.UserStatus, error) {
	client, conn, err := apiClient(ip)
	if err != nil {
		return nil, err
	}
	stream, err := client.ListUsers(t.ctx, &service.ListUsersRequest{})
	if err != nil {
		logrus.Errorf("%s list users stream err: %v\n", ip, err)
		return nil, errors.New(constant.GrpcError)
	}
	defer func() {
		stream.CloseSend()
		conn.Close()
	}()

	var userStatus []*service.UserStatus
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			logrus.Errorf("%s list users recv err: %v\n", ip, err)
		}
		userStatus = append(userStatus, resp.Status)
	}
	return userStatus, nil
}

// 查询节点上的用户
func (t *trojanGoApi) GetUser(ip string, hash string) (*service.UserStatus, error) {
	client, conn, err := apiClient(ip)
	if err != nil {
		return nil, err
	}
	stream, err := client.GetUsers(t.ctx)
	if err != nil {
		logrus.Errorf("%s get users stream err: %v\n", ip, err)
		return nil, errors.New(constant.GrpcError)
	}
	defer func() {
		stream.CloseSend()
		conn.Close()
	}()

	err = stream.Send(&service.GetUsersRequest{
		User: &service.User{
			Hash: hash,
		},
	})
	if err != nil {
		logrus.Errorf("%s get users stream send err: %v\n", ip, err)
		return nil, errors.New(constant.GrpcError)
	}
	resp, err := stream.Recv()
	if err != nil {
		logrus.Errorf("%s get users stream recv err: %v\n", ip, err)
		return nil, errors.New(constant.GrpcError)
	}
	return resp.Status, nil
}

// 节点上设置用户
func (t *trojanGoApi) setUser(ip string, setUsersRequest *service.SetUsersRequest) error {
	client, conn, err := apiClient(ip)
	if err != nil {
		return err
	}
	stream, err := client.SetUsers(t.ctx)
	if err != nil {
		logrus.Errorf("%s set users stream err: %v\n", ip, err)
		return errors.New(constant.GrpcError)
	}
	defer func() {
		stream.CloseSend()
		conn.Close()
	}()

	err = stream.Send(setUsersRequest)
	if err != nil {
		logrus.Errorf("%s set user stream send err: %v\n", ip, err)
		return errors.New(constant.GrpcError)
	}
	resp, err := stream.Recv()
	if err != nil {
		logrus.Errorf("%s set user stream recv err: %v\n", ip, err)
		return errors.New(constant.GrpcError)
	}
	if !resp.Success {
		logrus.Errorf("%s set user fail err: %v\n", ip, err)
		// 重试
	}
	return nil
}

// 节点上设置用户设备数
func (t *trojanGoApi) SetUserIpLimit(ip string, hash string, ipLimit int) error {
	req := &service.SetUsersRequest{
		Status: &service.UserStatus{
			User: &service.User{
				Hash: hash,
			},
			IpLimit: int32(ipLimit),
		},
		Operation: service.SetUsersRequest_Modify,
	}
	return t.setUser(ip, req)
}

// 节点上设置用户限速
func (t *trojanGoApi) SetUserSpeedLimit(ip string, hash string, uploadSpeedLimit int, downloadSpeedLimit int) error {
	req := &service.SetUsersRequest{
		Status: &service.UserStatus{
			User: &service.User{
				Hash: hash,
			},
			SpeedLimit: &service.Speed{
				UploadSpeed:   uint64(uploadSpeedLimit),
				DownloadSpeed: uint64(downloadSpeedLimit),
			},
		},
		Operation: service.SetUsersRequest_Modify,
	}
	return t.setUser(ip, req)
}

// 节点上删除用户
func (t *trojanGoApi) DeleteUser(ip string, hash string) error {
	req := &service.SetUsersRequest{
		Status: &service.UserStatus{
			User: &service.User{
				Hash: hash,
			},
		},
		Operation: service.SetUsersRequest_Delete,
	}
	return t.setUser(ip, req)
}

// 节点上添加用户
func (t *trojanGoApi) AddUser(ip string, hash string, ipLimit int, uploadSpeedLimit int, downloadSpeedLimit int) error {
	req := &service.SetUsersRequest{
		Status: &service.UserStatus{
			User: &service.User{
				Hash: hash,
			},
			IpLimit: int32(ipLimit),
			SpeedLimit: &service.Speed{
				UploadSpeed:   uint64(uploadSpeedLimit),
				DownloadSpeed: uint64(downloadSpeedLimit),
			},
		},
		Operation: service.SetUsersRequest_Add,
	}
	return t.setUser(ip, req)
}
