package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestMd5(t *testing.T) {
	s := StringToMd5("s")
	fmt.Println(s)
}

func TestInterfaces(t *testing.T) {
	addr := GetLocalIP()
	fmt.Println(addr)
}

func TestIPToDec(t *testing.T) {
	now := time.Now().Unix()
	fmt.Println(now)
	addr := "10.46.212.144"
	ipNum, _ := IPToDec(addr)
	fmt.Println(ipNum)
}
