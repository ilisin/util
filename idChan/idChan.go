package idChan

import (
	"gogs.xlh/tools/util/logger"
	"gogs.xlh/tools/util/tools"
)

var IDChan chan string
var ExitChan chan bool

var bInit bool

func IdDump(bufferLen int) {
	IDChan = make(chan string, bufferLen)
	ExitChan = make(chan bool)
	if bInit {
		return
	}
	bInit = true
	for {
		id := tools.Guid()
		select {
		case IDChan <- id:
		case <-ExitChan:
			goto exit
		}
	}
	bInit = false
exit:
	logger.Info("ID: closing")
}
