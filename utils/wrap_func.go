package utils

import (
	"fmt"
	"runtime"
	"strings"
)

const UnKnown = "UN_KNOWN"

type Wrapper struct {
	fileName   string // 文件名
	funcName   string // 方法名
	lineNumber int    // 行数
}

func NewWrapper() *Wrapper {
	fileName := UnKnown
	funcName := UnKnown
	fileLine := 0
	if pc, codePath, line, ok := runtime.Caller(1); ok {
		fileName = getLastName(codePath)
		funcName = getLastName(runtime.FuncForPC(pc).Name())
		fileLine = line
	}
	return &Wrapper{
		fileName:   fileName,
		funcName:   funcName,
		lineNumber: fileLine,
	}
}

func (w *Wrapper) String() string {
	return fmt.Sprintf("[%s:%s:%d]", w.fileName, w.funcName, w.lineNumber)
}

// name = "aa/bb" res = "bb"
func getLastName(name string) (res string) {
	return name[strings.LastIndex(name, "/")+1:]
}
