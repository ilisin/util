package crontab

import (
	"gogs.xlh/tools/util/crontab/job"
	"gogs.xlh/tools/util/crontab/provider"
)

func Run() {
	s_engine = Engine{}
	s_engine.Provider = provider.Provider{}
	s_engine.Provider.OutStream = make(chan job.BasicJob, 128)
	s_engine.Start()
}
