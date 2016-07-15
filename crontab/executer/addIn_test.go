package executer

import (
	"gogs.xlh/tools/util/logger"
	"testing"
)

func TestAddIn(t *testing.T) {
	addIn := &AddIn{}
	rc, err := addIn.Exec("/home/gaoguangting/go/src/imooly.com/hello")
	if err != nil {
		t.Error("回调错误")
	}
	<-rc
	logger.Trace("执行命令OK")
}
