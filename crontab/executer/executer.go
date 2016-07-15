package executer

import "gogs.xlh/tools/util/logger"

type Executer interface {
	Exec(cmd string, params ...string) (chan bool, error)
}

func init() {
	// register the system job
	RegisterSystemJob("JOB_SYS_0001", Demo)
}

func Demo() error {
	logger.Trace("DEMO注册函数")
	return nil
}
