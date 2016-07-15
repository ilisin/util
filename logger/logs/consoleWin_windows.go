/**
在win32环境中，控制台的文字和背景可以通过动态链接库kernel32.dll的一个函数SetConsoleTextAttribute()这个函数实现。 这个函数接受一个标准输出的handle作为第一个参数，第二个参数是用来控制颜色的attribute。控制台的颜色分成16种不同的值。 每个都可以用一个16进制的数来表示。

分别是：

0 = Black
1 = Blue
2 = Green
3 = Aqua
4 = Red
5 = Purple
6 = Yellow
7 = White
8 = Gray
9 = Light Blue
A = Light Green
B = Light Aqua
C = Light Red
D = Light Purple
E = Light Yellow
F = Bright White
32位的高位表示背景，低位表示文字颜色。
*/
package logs

import (
	"fmt"
	"strings"
	"syscall"
)

const (
	//标准输出宏
	STD_OUTPUT_HANDLE = uint32(-11 & 0xFFFFFFFF)
)

/**
NewBrush("1;36"), // Trace      cyan
NewBrush("1;34"), // Debug      blue
NewBrush("1;32"), // Info       green
NewBrush("1;33"), // Warn       yellow
NewBrush("1;31"), // Error      red
NewBrush("1;35"), // Critical   purple
*/
const (
	LOG_TRACE = iota
	LOG_DEBUG
	LOG_INFO
	LOG_WARN
	LOG_ERROR
	LOG_CRITICAL
	LOG_UNKOWN
)

const (
	WINCON_BLACK       = 0x0
	WINCON_BLUE        = 0x1
	WINCON_GREEN       = 0x2
	WINCON_AQUA        = 0x3
	WINCON_RED         = 0x4
	WINCON_PURPLE      = 0x5
	WINCON_YELLOW      = 0x6
	WINCON_WHITE       = 0x7 //Unkown
	WINCON_GRAY        = 0x8
	WINCON_LIGHTBLUE   = 0x9 //Debug
	WINCON_LIGHTGREEN  = 0xa //Info
	WINCON_LIGHTAQUA   = 0xb //Trace
	WINCON_LIGHTRED    = 0xc //Error
	WINCON_LIGHTPURPLE = 0xd //Critical
	WINCON_LIGHTYELLOW = 0xe //Warn
	WINCON_LIGHTWHITE  = 0xf
)

type LogLevel int

var logColorMap = map[LogLevel]uint32{
	LOG_TRACE:    WINCON_LIGHTAQUA,
	LOG_DEBUG:    WINCON_LIGHTBLUE,
	LOG_INFO:     WINCON_LIGHTGREEN,
	LOG_WARN:     WINCON_LIGHTYELLOW,
	LOG_ERROR:    WINCON_LIGHTRED,
	LOG_CRITICAL: WINCON_LIGHTPURPLE,
	LOG_UNKOWN:   WINCON_LIGHTWHITE,
}

var (
	err         error
	kernel32, _ = syscall.LoadLibrary("kernel32.dll")
	//设置console属性
	setConsoleTextAttribute, _ = syscall.GetProcAddress(kernel32, "SetConsoleTextAttribute")
	//获取标准输入输出的函数
	getStdHandle, _ = syscall.GetProcAddress(kernel32, "GetStdHandle")
	//标准输出
	hCon uintptr
)

func init() {
	//nargs 代表参数个数
	var nargs uintptr = 1
	//参数需要全部转成uinptr
	hCon, _, _ = syscall.Syscall(uintptr(getStdHandle), nargs, uintptr(STD_OUTPUT_HANDLE), 0, 0)
}
func SetConsoleTextAttribute(hConsoleOutput uintptr, wAttributes uint32) bool {
	var nargs uintptr = 2
	ret, _, _ := syscall.Syscall(setConsoleTextAttribute, nargs, hConsoleOutput, uintptr(wAttributes), 0)
	return ret != 0
}

type formatNode struct {
	Level LogLevel
	Text  string
}

func ConsoleWinOut(level int, text string) {
	SetConsoleTextAttribute(hCon, logColorMap[LogLevel(level)])
	fmt.Println(text)
	SetConsoleTextAttribute(hCon, logColorMap[LOG_UNKOWN])
}

func ConsoleOutWithLinuxFmt(text string) {
	defer syscall.FreeLibrary(kernel32)
	arr := make([]formatNode, 0)
	for {
		if len(text) == 0 {
			break
		}
		i := strings.Index(text, BRUSH_PRE_LINUX)
		if i < 0 {
			break
		}
		node := formatNode{}
		e := i
		if i > 0 {
			node.Level = LOG_UNKOWN
			node.Text = text[0:i]
			text = text[i:]
		} else {
			e = strings.Index(text, BRUSH_RESET_LINUX)
			temp := strings.TrimLeft(text, BRUSH_PRE_LINUX)
			s := len(BRUSH_PRE_LINUX) + len("1;36")
			node.Text = temp[s:e]
			if strings.HasPrefix(temp, "1;36") {
				node.Level = LOG_TRACE
			} else if strings.HasPrefix(temp, "1;34") {
				node.Level = LOG_DEBUG
			} else if strings.HasPrefix(temp, "1;32") {
				node.Level = LOG_INFO
			} else if strings.HasPrefix(temp, "1;33") {
				node.Level = LOG_WARN
			} else if strings.HasPrefix(temp, "1;31") {
				node.Level = LOG_ERROR
			} else if strings.HasPrefix(temp, "1;35") {
				node.Level = LOG_CRITICAL
			}
			e += len(BRUSH_RESET_LINUX)
		}
		arr = append(arr, node)
		text = text[e:]
	}
	if len(text) > 0 {
		node := formatNode{}
		node.Level = LOG_UNKOWN
		node.Text = text
		arr = append(arr, node)
	}
	for _, n := range arr {
		SetConsoleTextAttribute(hCon, logColorMap[n.Level])
		fmt.Println(n.Text)
		SetConsoleTextAttribute(hCon, logColorMap[LOG_UNKOWN])
	}
}
