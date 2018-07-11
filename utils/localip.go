package utils

import (
	"fmt"
	"strings"
)

var localIP string
var localPort string
var localCallback string

func SetLocalCallback(listenAddr string) {
	addrs := GetLocalIPs()
	for _, address := range addrs {
		if strings.HasPrefix(address, "10.") == true {
			localIP = address
			break
		}
	}
	if len(localIP) == 0 {
		localIP = GetLocalIP()
	}

	pos := strings.Index(listenAddr, ":")
	if pos != -1 {
		localPort = listenAddr[pos+1:]
	} else {
		localPort = "80"
	}
	if len(localIP) > 0 {
		localCallback = fmt.Sprintf("http://%s:%s", localIP, localPort)
	}
	return
}

func GetLocalCallback() string {
	return localCallback
}
