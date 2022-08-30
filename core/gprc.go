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
	"trojan/module/constant"
)

func newGrpcInstance(ip string) (conn *grpc.ClientConn, ctx context.Context, clo func(), err error) {
	conn, err = grpc.Dial(fmt.Sprintf("%s:%d", ip, 8100),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
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
func AddNode(ip string, nodeAddDto *NodeAddDto) error {
	conn, ctx, clo, err := newGrpcInstance(ip)
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

func RemoveNode(ip string, nodeRemoveDto *NodeRemoveDto) error {
	conn, ctx, clo, err := newGrpcInstance(ip)
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
