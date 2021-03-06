package job

import (
	"gogs.xlh/tools/util/logger"
	"testing"
	"time"
)

func TestBasicJob(t *testing.T) {
	bj := BasicJob{
		Id:       "J001",
		ExeTime:  time.Now().Add(5 * time.Second),
		Type:     JOB_SYSTEM,
		Live:     5,
		Interval: 3 * time.Second,
		//Command : "http://192.168.0.108:8080/v1/user",
		Command: "JOB_SYS_0001",
		Params:  ""}
	logger.Info("同步自行开始")
	if err := bj.StartDo(); err != nil {
		t.Fatal("执行任务异常")
	}
	logger.Info("同步执行OK")

	logger.Info("异步开始执行")
	bj2 := BasicJob{
		Id:       "J002",
		ExeTime:  time.Now().Add(6 * time.Second),
		Type:     JOB_SYSTEM,
		Live:     3,
		Interval: 3 * time.Second,
		Command:  "JOB_SYS_0001",
		Params:   ""}
	go bj2.StartDo()
	logger.Info("执行OK")
	time.Sleep(20 * time.Second)
}
