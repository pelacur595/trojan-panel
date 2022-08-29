package util

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	encode, err := AesEncode("123456")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("加密后为：%s\n", encode)
	decode, err := AesDecode(encode)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("解密后为：%s\n", decode)
}
