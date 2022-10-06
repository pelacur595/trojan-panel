package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
	"trojan-panel/module/constant"
)

func newGrpcInstance(token string, ip string) (conn *grpc.ClientConn, ctx context.Context, clo func(), err error) {
	tokenParam := TokenValidateParam{
		Token: token,
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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

func AddNode(token string, ip string, nodeAddDto *NodeAddDto) error {
	conn, ctx, clo, err := newGrpcInstance(token, ip)
	defer clo()
	if err != nil {
		return err
	}
	client := NewApiNodeServiceClient(conn)
	send, err := client.AddNode(ctx, nodeAddDto)
	if err != nil {
		log.Println(err)
		return errors.New(constant.GrpcError)
	}
	if send.Success {
		return nil
	}
	return errors.New(send.Msg)
}

func RemoveNode(token string, ip string, nodeRemoveDto *NodeRemoveDto) error {
	conn, ctx, clo, err := newGrpcInstance(token, ip)
	defer clo()
	if err != nil {
		return err
	}
	client := NewApiNodeServiceClient(conn)
	send, err := client.RemoveNode(ctx, nodeRemoveDto)
	if err != nil {
		log.Println(err)
		return errors.New(constant.GrpcError)
	}
	if send.Success {
		return nil
	}
	return errors.New(send.Msg)
}

func RemoveAccount(token string, ip string, accountRemoveDto *AccountRemoveDto) error {
	conn, ctx, clo, err := newGrpcInstance(token, ip)
	defer clo()
	if err != nil {
		return err
	}
	client := NewApiAccountServiceClient(conn)
	send, err := client.RemoveAccount(ctx, accountRemoveDto)
	if err != nil {
		log.Println(err)
		return errors.New(constant.GrpcError)
	}
	if send.Success {
		return nil
	}
	return errors.New(send.Msg)
}
