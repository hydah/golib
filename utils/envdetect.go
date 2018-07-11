package utils

import (
	"flag"
	"strings"
)

func IsForTestCase() bool {
	/*
	 * 当前环境是否是在测试环境中
	 */
	if flag.Lookup("test.v") == nil {
		return false
	}
	return true
}

func GetTestArg(argName string) (argVal string) {
	/*
	 * 获取 go test 后面 -args 里面的参数。
	 * 例如 go test -args aaa=bbb, 使用时用 argName=aaa 即可
	 */
	argVal = ""
	for _, arg := range flag.Args() {
		arglist := strings.Split(arg, "=")
		if argName == arglist[0] {
			argVal = arglist[1]
		}
	}
	return
}
