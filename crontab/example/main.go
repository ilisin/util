package main

import (
	"gogs.xlh/tools/util/crontab"
	"gogs.xlh/tools/util/logger"
	"time"
)

var (
	logConfig string = `{
	"filename":"logs.log",
	"maxlines":10000,
	"maxsize":10000000,
	"daily":true,
	"maxdays":15,
	"rotate":true
	}`
)

func main() {
	logger.SetLogger("file", logConfig)
	crontab.Run()

	time.Sleep(1 * time.Hour)
}
