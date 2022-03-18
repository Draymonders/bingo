package log

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	Log *logrus.Logger
)

func Init(logLevel logrus.Level) {
	Log = logrus.New()
	Log.SetReportCaller(true)
	Log.SetLevel(logLevel)
	log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			//处理文件名
			fileName := path.Base(frame.File)
			//fmt.Printf("filename=%v,function=%v line=%v\n", fileName, frame.Function, frame.Line)
			f := frame.Function
			ind := strings.LastIndex(f, "/")
			funcName := f[ind+1:]
			return funcName, fmt.Sprintf("%v:%v", fileName, frame.Line)
		},
	})
}