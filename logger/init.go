package logger

import (
	"fmt"
	"gogs.xlh/tools/util/logger/logs"
)

func init() {

	// init Logger
	Logger = logs.NewLogger(1000)
	Logger.EnableFuncCallDepth(true)
	Logger.SetLogFuncCallDepth(3)
	err := Logger.SetLogger("console", "")
	if err != nil {
		fmt.Println("init console log error:", err)
	}
	SetLevel(LevelTrace) //设置日志级别
}
