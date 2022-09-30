package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"time"
	"trojan/module/constant"
)

func newGrpcInstance(ip string, token string) (conn *grpc.ClientConn, ctx context.Context, clo func(), err error) {
	tokenParam := TokenValidateParam{
		Token: token,
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(&tokenParam),
	}
	conn, err = grpc.Dial(fmt.Sprintf("%s:%d", ip, 8100),
		opts...)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	clo = func() {
		cancel()
		conn.Close()
	}
	if err != nil {
		logrus.Errorf("gRPC初始化失败 err: %v\n", err)
		err = errors.New(constant.GrpcError)
	}
	return
}
func AddNode(ip string, token string, nodeAddDto *NodeAddDto) error {
	conn, ctx, clo, err := newGrpcInstance(ip, token)
	defer clo()
	if err != nil {
		return err
	}
	client := NewApiNodeServiceClient(conn)
	send, err := client.AddNode(ctx, nodeAddDto)
	if err != nil {
		log.Println(err)
	}
	if send.Success {
		return nil
	}
	return errors.New(send.Msg)
}

func RemoveNode(ip string, token string, nodeRemoveDto *NodeRemoveDto) error {
	conn, ctx, clo, err := newGrpcInstance(ip, token)
	defer clo()
	if err != nil {
		return err
	}
	client := NewApiNodeServiceClient(conn)
	send, err := client.RemoveNode(ctx, nodeRemoveDto)
	if err != nil {
		log.Println(err)
	}
	if send.Success {
		return nil
	}
	return errors.New(send.Msg)
}

func RemoveAccount(ip string, token string, accountRemoveDto *AccountRemoveDto) error {
	conn, ctx, clo, err := newGrpcInstance(ip, token)
	defer clo()
	if err != nil {
		return err
	}
	client := NewApiAccountServiceClient(conn)
	send, err := client.RemoveAccount(ctx, accountRemoveDto)
	if err != nil {
		log.Println(err)
	}
	if send.Success {
		return nil
	}
	return errors.New(send.Msg)
}
