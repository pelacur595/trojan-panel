package core

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	state, err := NodeServerState("", "127.0.0.1")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(state.CpuUsed)
		fmt.Println(state.MemUsed)
		fmt.Println(state.DiskUsed)
	}
}
