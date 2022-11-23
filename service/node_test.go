package service

import (
	"fmt"
	"testing"
	"trojan-panel/core"
)

func TestGrpcAddNode(t *testing.T) {
	dto := core.NodeAddDto{
		NodeTypeId: 4,
		Port:       4444,
		Ip:         "demo.wellveryfunny.xyz",
	}
	if err := core.AddNode("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50Vm8iOnsiaWQiOjEsInF1b3RhIjowLCJkb3dubG9hZCI6MCwidXBsb2FkIjowLCJ1c2VybmFtZSI6InN5c2FkbWluIiwiZW1haWwiOiIiLCJyb2xlSWQiOjEsImRlbGV0ZWQiOjAsImV4cGlyZVRpbWUiOjAsImNyZWF0ZVRpbWUiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiJ9LCJleHAiOjE2NjQwMTQyNjksImlzcyI6InRyb2phbi1wYW5lbCJ9.vUbGGp42XTyndNNH01aSj6YW_bfck-jmzUcs1JtVMb0",
		"127.0.0.1", &dto); err != nil {
		fmt.Println(err.Error())
	}
}

func TestGrpcRemoveNode(t *testing.T) {
	removeDto := core.NodeRemoveDto{NodeType: 2, Port: 443}
	if err := core.RemoveNode("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50Vm8iOnsiaWQiOjEsInF1b3RhIjowLCJkb3dubG9hZCI6MCwidXBsb2FkIjowLCJ1c2VybmFtZSI6InN5c2FkbWluIiwiZW1haWwiOiIiLCJyb2xlSWQiOjEsImRlbGV0ZWQiOjAsImV4cGlyZVRpbWUiOjAsImNyZWF0ZVRpbWUiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiJ9LCJleHAiOjE2NjQwMTQyNjksImlzcyI6InRyb2phbi1wYW5lbCJ9.vUbGGp42XTyndNNH01aSj6YW_bfck-jmzUcs1JtVMb0",
		"127.0.0.1", &removeDto); err != nil {
		fmt.Println(err.Error())
	}
}
