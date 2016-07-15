package logs

import (
	"fmt"
)

func ConsoleWinOut(level int, text string) {
	fmt.Println(text)
}

func ConsoleOutWithLinuxFmt(text string) {
	fmt.Println(text)
}
