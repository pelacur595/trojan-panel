package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"time"
	"trojan-panel/module/constant"
)

func newGrpcInstance(token string, ip string, grpcPort uint, timeout time.Duration) (conn *grpc.ClientConn, ctx context.Context, clo func(), err error) {
	tokenParam := TokenValidateParam{
		Token: token,
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&tokenParam),
	}
	conn, err = grpc.Dial(fmt.Sprintf("%s:%d", ip, grpcPort),
		opts...)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	clo = func() {
		cancel()
		conn.Close()
	}
	if err != nil {
		logrus.Errorf("gRPC初始化失败 err: %v", err)
		err = errors.New(constant.GrpcError)
	}
	return
}

func AddNode(token string, ip string, grpcPort uint, nodeAddDto *NodeAddDto) error {
	conn, ctx, clo, err := newGrpcInstance(token, ip, grpcPort, 4*time.Second)
	defer clo()
	if err != nil {
		return err
	}
	client := NewApiNodeServiceClient(conn)
	send, err := client.AddNode(ctx, nodeAddDto)
	if err != nil {
		logrus.Errorf("gRPC添加节点异常 ip: %s grpc port: %d err: %v", ip, grpcPort, err)
		return errors.New(constant.GrpcError)
	}
	if send.Success {
		return nil
	}
	return errors.New(send.Msg)
}

func RemoveNode(token string, ip string, grpcPort uint, nodeRemoveDto *NodeRemoveDto) error {
	conn, ctx, clo, err := newGrpcInstance(token, ip, grpcPort, 4*time.Second)
	defer clo()
	if err != nil {
		return err
	}
	client := NewApiNodeServiceClient(conn)
	send, err := client.RemoveNode(ctx, nodeRemoveDto)
	if err != nil {
		logrus.Errorf("gRPC删除节点异常 ip: %s grpc port: %d err: %v", ip, grpcPort, err)
		return errors.New(constant.GrpcError)
	}
	if send.Success {
		return nil
	}
	return errors.New(send.Msg)
}

func RemoveAccount(token string, ip string, grpcPort uint, accountRemoveDto *AccountRemoveDto) error {
	conn, ctx, clo, err := newGrpcInstance(token, ip, grpcPort, 4*time.Second)
	defer clo()
	if err != nil {
		return err
	}
	client := NewApiAccountServiceClient(conn)
	send, err := client.RemoveAccount(ctx, accountRemoveDto)
	if err != nil {
		logrus.Errorf("gRPC删除用户异常 ip: %s grpc porr: %d err: %v", ip, grpcPort, err)
		return errors.New(constant.GrpcError)
	}
	if send.Success {
		return nil
	}
	return errors.New(send.Msg)
}

// GetNodeState 查询节点状态
func GetNodeState(token string, ip string, grpcPort uint) (bool, error) {
	conn, ctx, clo, err := newGrpcInstance(token, ip, grpcPort, 2*time.Second)
	defer clo()
	if err != nil {
		return false, err
	}
	client := NewApiStateServiceClient(conn)
	nodeStateDto := NodeStateDto{}
	send, err := client.GetNodeState(ctx, &nodeStateDto)
	if err != nil {
		logrus.Errorf("gRPC GetNodeState 异常 ip: %s grpc port: %d err: %v", ip, grpcPort, err)
		return false, errors.New(constant.GrpcError)
	}
	if send.Success {
		return true, nil
	}
	logrus.Errorf("gRPC GetNodeState 失败 ip: %s grpc port: %d err: %v", ip, grpcPort, err)
	return false, errors.New(send.Msg)
}

// GetNodeServerState 查询服务器状态
func GetNodeServerState(token string, ip string, grpcPort uint) (*NodeServerStateVo, error) {
	conn, ctx, clo, err := newGrpcInstance(token, ip, grpcPort, 2*time.Second)
	defer clo()
	if err != nil {
		return nil, err
	}
	client := NewApiStateServiceClient(conn)
	nodeServerStateDto := NodeServerStateDto{}
	send, err := client.GetNodeServerState(ctx, &nodeServerStateDto)
	if err != nil {
		logrus.Errorf("gRPC GetNodeServerState 异常 ip: %s grpc port: %d err: %v", ip, grpcPort, err)
		return nil, errors.New(constant.GrpcError)
	}
	if send.Success {
		var nodeServerStateVo NodeServerStateVo
		if err = anypb.UnmarshalTo(send.Data, &nodeServerStateVo, proto.UnmarshalOptions{}); err != nil {
			logrus.Errorf("gRPC GetNodeServerState 返序列化异常 ip: %s grpc port: %d err: %v", ip, grpcPort, err)
			return nil, errors.New(constant.GrpcError)
		}
		return &nodeServerStateVo, nil
	}
	logrus.Errorf("gRPC GetNodeServerState 失败 ip: %s grpc port: %d err: %v", ip, grpcPort, err)
	return nil, errors.New(send.Msg)
}

func GetNodeServerInfo(token string, ip string, grpcPort uint) (*NodeServerInfoVo, error) {
	conn, ctx, clo, err := newGrpcInstance(token, ip, grpcPort, 2*time.Second)
	defer clo()
	if err != nil {
		return nil, err
	}
	client := NewApiNodeServerServiceClient(conn)
	nodeServerInfoDto := NodeServerInfoDto{}
	send, err := client.GetNodeServerInfo(ctx, &nodeServerInfoDto)
	if err != nil {
		logrus.Errorf("gRPC 查询服务器状态 异常 ip: %s grpc port: %d err: %v", ip, grpcPort, err)
		return nil, errors.New(constant.GrpcError)
	}
	if send.Success {
		var nodeServerInfoVo NodeServerInfoVo
		if err = anypb.UnmarshalTo(send.Data, &nodeServerInfoVo, proto.UnmarshalOptions{}); err != nil {
			logrus.Errorf("查询服务器状态返序列化异常 ip: %s grpc port: %d err: %v", ip, grpcPort, err)
			return nil, errors.New(constant.GrpcError)
		}
		return &nodeServerInfoVo, nil
	}
	return nil, errors.New(send.Msg)
}
