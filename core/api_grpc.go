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

func newGrpcInstance(token string, ip string, timeout time.Duration) (conn *grpc.ClientConn, ctx context.Context, clo func(), err error) {
	tokenParam := TokenValidateParam{
		Token: token,
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&tokenParam),
	}
	conn, err = grpc.Dial(fmt.Sprintf("%s:%d", ip, 8100),
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

func AddNode(token string, ip string, nodeAddDto *NodeAddDto) error {
	conn, ctx, clo, err := newGrpcInstance(token, ip, 4*time.Second)
	defer clo()
	if err != nil {
		return err
	}
	client := NewApiNodeServiceClient(conn)
	send, err := client.AddNode(ctx, nodeAddDto)
	if err != nil {
		logrus.Errorf("gRPC添加节点异常 ip: %s err: %v", ip, err)
		return errors.New(constant.GrpcError)
	}
	if send.Success {
		return nil
	}
	return errors.New(send.Msg)
}

func RemoveNode(token string, ip string, nodeRemoveDto *NodeRemoveDto) error {
	conn, ctx, clo, err := newGrpcInstance(token, ip, 4*time.Second)
	defer clo()
	if err != nil {
		return err
	}
	client := NewApiNodeServiceClient(conn)
	send, err := client.RemoveNode(ctx, nodeRemoveDto)
	if err != nil {
		logrus.Errorf("gRPC删除节点异常 ip: %s err: %v", ip, err)
		return errors.New(constant.GrpcError)
	}
	if send.Success {
		return nil
	}
	return errors.New(send.Msg)
}

func RemoveAccount(token string, ip string, accountRemoveDto *AccountRemoveDto) error {
	conn, ctx, clo, err := newGrpcInstance(token, ip, 4*time.Second)
	defer clo()
	if err != nil {
		return err
	}
	client := NewApiAccountServiceClient(conn)
	send, err := client.RemoveAccount(ctx, accountRemoveDto)
	if err != nil {
		logrus.Errorf("gRPC删除用户异常 ip: %s err: %v", ip, err)
		return errors.New(constant.GrpcError)
	}
	if send.Success {
		return nil
	}
	return errors.New(send.Msg)
}

func Ping(token string, ip string) (bool, error) {
	conn, ctx, clo, err := newGrpcInstance(token, ip, time.Second)
	defer clo()
	if err != nil {
		return false, err
	}
	client := NewApiStateServiceClient(conn)
	stateDto := StateDto{}
	send, err := client.Ping(ctx, &stateDto)
	if err != nil {
		logrus.Errorf("gRPC ping 异常 ip: %s err: %v", ip, err)
		return false, errors.New(constant.GrpcError)
	}
	if send.Success {
		return true, nil
	}
	return false, errors.New(send.Msg)
}

func NodeServerState(token string, ip string) (*NodeServerGroupVo, error) {
	conn, ctx, clo, err := newGrpcInstance(token, ip, 2*time.Second)
	defer clo()
	if err != nil {
		return nil, err
	}
	client := NewApiNodeServerServiceClient(conn)
	nodeServerGroupDto := NodeServerGroupDto{}
	send, err := client.NodeServerState(ctx, &nodeServerGroupDto)
	if err != nil {
		logrus.Errorf("gRPC 查询服务器状态 异常 ip: %s err: %v", ip, err)
		return nil, errors.New(constant.GrpcError)
	}
	if send.Success {
		var nodeServerGroupVo NodeServerGroupVo
		if err = anypb.UnmarshalTo(send.Data, &nodeServerGroupVo, proto.UnmarshalOptions{}); err != nil {
			logrus.Errorf("查询服务器状态返序列化异常 ip: %s err: %v", ip, err)
			return nil, errors.New(constant.GrpcError)
		}
		return &nodeServerGroupVo, nil
	}
	return nil, errors.New(send.Msg)
}
