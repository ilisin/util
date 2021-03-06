package logs

import (
	"encoding/json"
	"log"
	"os"
	"runtime"
)

const (
	BRUSH_PRE_LINUX   = "\033["
	BRUSH_RESET_LINUX = "\033[0m"
)

type Brush func(string) string

func NewBrush(color string) Brush {
	return func(text string) string {
		return BRUSH_PRE_LINUX + color + "m" + text + BRUSH_RESET_LINUX
	}
}

var colors = []Brush{
	NewBrush("1;36"), // Trace      cyan
	NewBrush("1;34"), // Debug      blue
	NewBrush("1;32"), // Info       green
	NewBrush("1;33"), // Warn       yellow
	NewBrush("1;31"), // Error      red
	NewBrush("1;35"), // Critical   purple
}

// ConsoleWriter implements LoggerInterface and writes messages to terminal.
type ConsoleWriter struct {
	lg    *log.Logger
	Level int `json:"level"`
}

// create ConsoleWriter returning as LoggerInterface.
func NewConsole() LoggerInterface {
	cw := new(ConsoleWriter)
	cw.lg = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	cw.Level = LevelTrace
	return cw
}

// init console logger.
// jsonconfig like '{"level":LevelTrace}'.
func (c *ConsoleWriter) Init(jsonconfig string) error {
	if len(jsonconfig) == 0 {
		return nil
	}
	err := json.Unmarshal([]byte(jsonconfig), c)
	if err != nil {
		return err
	}
	return nil
}

// write message in console.
func (c *ConsoleWriter) WriteMsg(msg string, level int) error {
	if level < c.Level {
		return nil
	}
	if goos := runtime.GOOS; goos == "windows" {
		//c.lg.Println(msg)
		ConsoleWinOut(level, msg)
	} else {
		c.lg.Println(colors[level](msg))
	}
	return nil
}

// implementing method. empty.
func (c *ConsoleWriter) Destroy() {

}

// implementing method. empty.
func (c *ConsoleWriter) Flush() {

}

func init() {
	Register("console", NewConsole)
}
