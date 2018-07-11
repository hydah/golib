package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func IPToDec(ip string) (ipNum int64, err error) {
	var dec int64
	var hexStr string
	ipArr := strings.Split(ip, ".")
	for i := len(ipArr) - 1; i >= 0; i-- {
		dec, err = strconv.ParseInt(ipArr[i], 10, 64)
		if err != nil {
			return
		}
		if dec > 255 {
			err = fmt.Errorf("invail ip")
			return
		}
		hex := strconv.FormatInt(dec, 16)
		if len(hex) == 1 {
			hex = fmt.Sprintf("0%s", hex)
		}
		hexStr += hex
	}
	ipNum, err = strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		return
	}
	return
}

func DecToIP(ipNum int64) (ip string, err error) {
	hexStr := strconv.FormatInt(ipNum, 16)
	if len(hexStr) == 7 {
		hexStr = fmt.Sprintf("0%s", hexStr)
	}
	if len(hexStr) != 8 {
		err = fmt.Errorf("invail ipNum")
		return
	}
	var ipN int64
	hexArr := strings.Split(hexStr, "")
	ip = ""
	for i := 7; i >= 0; i -= 2 {
		ipN, err = strconv.ParseInt(fmt.Sprintf("%s%s", hexArr[i-1], hexArr[i]), 16, 64)
		if err != nil {
			return
		}
		ip = fmt.Sprintf("%s%d", ip, ipN)
		if i != 1 {
			ip = fmt.Sprintf("%s.", ip)
		}
	}
	return
}
