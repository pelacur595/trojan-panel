package service

import (
	"fmt"
	"testing"
	"trojan-panel/core"
)

func TestGrpcAddNode(t *testing.T) {
	dto := core.NodeAddDto{
		NodeTypeId: 4,
		Port:       883,
		Domain:     "demo.wellveryfunny.xyz",
	}
	if err := core.AddNode("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50Vm8iOnsiaWQiOjEsInF1b3RhIjowLCJkb3dubG9hZCI6MCwidXBsb2FkIjowLCJ1c2VybmFtZSI6InN5c2FkbWluIiwiZW1haWwiOiIiLCJyb2xlSWQiOjEsImRlbGV0ZWQiOjAsImV4cGlyZVRpbWUiOjAsImNyZWF0ZVRpbWUiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiIsInJvbGVzIjpbInN5c2FkbWluIiwiYWRtaW4iLCJ1c2VyIl19LCJleHAiOjE2NjkxOTgxMDAsImlzcyI6InRyb2phbi1wYW5lbCJ9.9DYvvo0F-XS55WUy5J8Lmj8gim243yDd_BiZZnNv0is",
		"127.0.0.1", 8100, &dto); err != nil {
		fmt.Println(err.Error())
	}
}

func TestGrpcRemoveNode(t *testing.T) {
	removeDto := core.NodeRemoveDto{NodeTypeId: 4, Port: 883}
	if err := core.RemoveNode("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50Vm8iOnsiaWQiOjEsInF1b3RhIjowLCJkb3dubG9hZCI6MCwidXBsb2FkIjowLCJ1c2VybmFtZSI6InN5c2FkbWluIiwiZW1haWwiOiIiLCJyb2xlSWQiOjEsImRlbGV0ZWQiOjAsImV4cGlyZVRpbWUiOjAsImNyZWF0ZVRpbWUiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiIsInJvbGVzIjpbInN5c2FkbWluIiwiYWRtaW4iLCJ1c2VyIl19LCJleHAiOjE2NjkxOTgxMDAsImlzcyI6InRyb2phbi1wYW5lbCJ9.9DYvvo0F-XS55WUy5J8Lmj8gim243yDd_BiZZnNv0is",
		"127.0.0.1", 8100, &removeDto); err != nil {
		fmt.Println(err.Error())
	}
}
