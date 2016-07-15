package logs

import (
	"fmt"
	"path"
	"runtime"
	"sync"
)

const (
	// log message levels
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelCritical
)

type loggerType func() LoggerInterface

// LoggerInterface defines the behavior of a log provider.
type LoggerInterface interface {
	Init(config string) error
	WriteMsg(msg string, level int) error
	Destroy()
	Flush()
}

var adapters = make(map[string]loggerType)

// Register makes a log provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, log loggerType) {
	if log == nil {
		panic("logs: Register provide is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("logs: Register called twice for provider " + name)
	}
	adapters[name] = log
}

// it can contain several providers and log message into all providers.
type Logger struct {
	lock                sync.Mutex
	level               int
	enableFuncCallDepth bool
	loggerFuncCallDepth int
	msg                 chan *logMsg
	outputs             map[string]LoggerInterface
}

type logMsg struct {
	level int
	msg   string
}

// NewLogger returns a new Logger.
// channellen means the number of messages in chan.
// if the buffering chan is full, logger adapters write to file or other way.
func NewLogger(channellen int64) *Logger {
	l := new(Logger)
	l.loggerFuncCallDepth = 2
	l.msg = make(chan *logMsg, channellen)
	l.outputs = make(map[string]LoggerInterface)
	//bl.SetLogger("console", "") // default output to console
	go l.startLogger()
	return l
}

// SetLogger provides a given logger adapter into Logger with config string.
// config need to be correct JSON as string: {"interval":360}.
func (l *Logger) SetLogger(adaptername string, config string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if log, ok := adapters[adaptername]; ok {
		lg := log()
		err := lg.Init(config)
		l.outputs[adaptername] = lg
		if err != nil {
			fmt.Println("logs.Logger.SetLogger: " + err.Error())
			return err
		}
	} else {
		return fmt.Errorf("logs: unknown adaptername %q (forgotten Register?)", adaptername)
	}
	return nil
}

// remove a logger adapter in Logger.
func (l *Logger) DelLogger(adaptername string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if lg, ok := l.outputs[adaptername]; ok {
		lg.Destroy()
		delete(l.outputs, adaptername)
		return nil
	} else {
		return fmt.Errorf("logs: unknown adaptername %q (forgotten Register?)", adaptername)
	}
}

func (l *Logger) writerMsg(loglevel int, msg string) error {
	if l.level > loglevel {
		return nil
	}
	lm := new(logMsg)
	lm.level = loglevel
	if l.enableFuncCallDepth {
		_, file, line, ok := runtime.Caller(l.loggerFuncCallDepth)
		if ok {
			_, filename := path.Split(file)
			lm.msg = fmt.Sprintf("[%s:%d] %s", filename, line, msg)
		} else {
			lm.msg = msg
		}
	} else {
		lm.msg = msg
	}
	l.msg <- lm
	return nil
}

// set log message level.
// if message level (such as LevelTrace) is less than logger level (such as LevelWarn), ignore message.
func (l *Logger) SetLevel(i int) {
	l.level = i
}

// set log funcCallDepth
func (l *Logger) SetLogFuncCallDepth(d int) {
	l.loggerFuncCallDepth = d
}

// enable log funcCallDepth
func (l *Logger) EnableFuncCallDepth(b bool) {
	l.enableFuncCallDepth = b
}

// start logger chan reading.
// when chan is full, write logs.
func (l *Logger) startLogger() {
	for {
		select {
		case bm := <-l.msg:
			for _, l := range l.outputs {
				l.WriteMsg(bm.msg, bm.level)
			}
		}
	}
}

// log trace level message.
func (l *Logger) Trace(format string, v ...interface{}) {
	msg := fmt.Sprintf("[T] "+format, v...)
	l.writerMsg(LevelTrace, msg)
}

// log debug level message.
func (l *Logger) Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf("[D] "+format, v...)
	l.writerMsg(LevelDebug, msg)
}

// log info level message.
func (l *Logger) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf("[I] "+format, v...)
	l.writerMsg(LevelInfo, msg)
}

// log warn level message.
func (l *Logger) Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf("[W] "+format, v...)
	l.writerMsg(LevelWarn, msg)
}

// log error level message.
func (l *Logger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf("[E] "+format, v...)
	l.writerMsg(LevelError, msg)
}

// log critical level message.
func (l *Logger) Critical(format string, v ...interface{}) {
	msg := fmt.Sprintf("[C] "+format, v...)
	l.writerMsg(LevelCritical, msg)
}

// flush all chan data.
func (l *Logger) Flush() {
	for _, lo := range l.outputs {
		lo.Flush()
	}
}

// close logger, flush all chan data and destroy all adapters in Logger.
func (l *Logger) Close() {
	for {
		if len(l.msg) > 0 {
			bm := <-l.msg
			for _, l := range l.outputs {
				l.WriteMsg(bm.msg, bm.level)
			}
		} else {
			break
		}
	}
	for _, l := range l.outputs {
		l.Flush()
		l.Destroy()
	}
}
