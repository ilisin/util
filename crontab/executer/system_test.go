package executer

import (
	"gogs.xlh/tools/util/logger"
	"testing"
)

func TestSystem(t *testing.T) {
	system := &System{}
	system.Exec("JOB_SYS_0001")
	logger.Trace("测试完成")
}
