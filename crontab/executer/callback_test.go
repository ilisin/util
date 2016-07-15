package executer

import (
	"gogs.xlh/tools/util/logger"
	"testing"
)

func TestCallback(t *testing.T) {
	callback := &Callback{}
	rc, err := callback.Exec("http://192.168.10.49:8080/v1/user")
	if err != nil {
		t.Error("回调错误")
	}
	<-rc
	logger.Trace("回调完成")
}
