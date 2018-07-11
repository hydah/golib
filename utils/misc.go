package utils

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func GetHostname() string {
	host, _ := os.Hostname()
	return host
}

func GetCurrPath() (string, string) {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	if index == -1 {
		return path, ""
	}
	p1 := path[:index]
	p2 := path[index:]
	return p1, p2
}

func GetFileExtName(path string) string {
	index := strings.LastIndex(path, string("."))
	if index == -1 {
		return ""
	}
	p := path[index+1:]
	return p
}

func HexDump(body []byte) {
	out := []byte(hex.Dump(body[:]))
	fmt.Printf("%s", out)
}

func HexDumpToString(body []byte) string {
	out := []byte(hex.Dump(body[:]))
	return fmt.Sprintf("%s", out)
}

func WritePidFile(prefix string) (err error) {
	filename := fmt.Sprintf("%s.pid", prefix)
	err = ioutil.WriteFile(filename, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
	return
}

func WritePid() (err error) {
	processName := os.Args[0]
	lastIndex := strings.LastIndex(processName, "/")
	processName = processName[lastIndex+1:]
	err = WritePidFile(processName)
	return
}

func RemovePid() (err error) {
	processName := os.Args[0]
	lastIndex := strings.LastIndex(processName, "/")
	processName = processName[lastIndex+1:]
	filename := fmt.Sprintf("%s.pid", processName)
	return os.Remove(filename)
}

func RedirectToPanicFile() {
	var discard *os.File
	var err error
	fileName := fmt.Sprintf(".panic_%d", os.Getpid())
	discard, err = os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		discard, err = os.OpenFile("/dev/null", os.O_RDWR, 0)
	}
	if err == nil {
		fd := discard.Fd()
		syscall.Dup2(int(fd), int(os.Stderr.Fd()))
	}
}

func InSlice(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func MatchInSlice(slice []string, item string) bool {
	for _, s := range slice {
		if strings.Index(s, item) != -1 {
			return true
		}
	}
	return false
}

func MergeIntSlice(src, des []int64) (ret []int64) {
	retMap := make(map[int64]bool)
	for _, id := range src {
		retMap[id] = true
	}

	for _, id := range des {
		retMap[id] = true
	}

	ret = make([]int64, 0)
	for id, _ := range retMap {
		ret = append(ret, id)
	}
	return
}
