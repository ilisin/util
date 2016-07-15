package tools

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func TimeParse(s string) (time.Time, error) {
	s = strings.Trim(s, " ")
	arr := strings.Split(s, " ")
	l := len(arr)
	if l == 0 {
		return time.Now(), errors.New("time string error" + fmt.Sprintf("l = %d", l))
	}

	if l == 1 {
		s += ` 00:00:00`
	}

	temp := "2006-01-02 15:04:05"
	if strings.Index(s, "/") > 0 {
		temp = "2006/1/2 15:04:05"
	}
	return time.ParseInLocation(temp, s, time.Local)
}

func TimeParseByDefault(s string, def time.Time) (r time.Time) {
	r = def
	s = strings.Trim(s, " ")
	arr := strings.Split(s, " ")
	l := len(arr)
	if l == 0 {
		return
	}

	if l == 1 {
		s += ` 00:00:00`
	}

	temp := "2006-01-02 15:04:05"
	if strings.Index(s, "/") > 0 {
		temp = "2006/1/2 15:04:05"
	}
	t, err := time.Parse(temp, s)
	if err == nil {
		r = t
	}
	return
}
