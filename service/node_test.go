package service

import (
	"fmt"
	"testing"
	"trojan/core"
)

func TestGrpcAddNode(t *testing.T) {
	dto := core.NodeAddDto{
		NodeTypeId:              2,
		TrojanGoPort:            443,
		TrojanGoIp:              "127.0.0.1",
		TrojanGoSni:             "",
		TrojanGoMuxEnable:       0,
		TrojanGoWebsocketEnable: 0,
		TrojanGoWebsocketPath:   "",
		TrojanGoWebsocketHost:   "",
		TrojanGoSSEnable:        0,
		TrojanGoSSMethod:        "",
		TrojanGoSSPassword:      "",
	}
	if err := core.AddNode("127.0.0.1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50Vm8iOnsiaWQiOjEsInF1b3RhIjowLCJkb3dubG9hZCI6MCwidXBsb2FkIjowLCJ1c2VybmFtZSI6InN5c2FkbWluIiwiZW1haWwiOiIiLCJyb2xlSWQiOjEsImRlbGV0ZWQiOjAsImV4cGlyZVRpbWUiOjAsImNyZWF0ZVRpbWUiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiJ9LCJleHAiOjE2NjQwMTQyNjksImlzcyI6InRyb2phbi1wYW5lbCJ9.vUbGGp42XTyndNNH01aSj6YW_bfck-jmzUcs1JtVMb0", &dto); err != nil {
		fmt.Println(err.Error())
	}
}
