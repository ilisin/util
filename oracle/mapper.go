package oracle

import (
	"strings"
)

type UpperMapper struct {
}

func (mapper UpperMapper) Obj2Table(name string) string {
	return strings.ToUpper(name)
}

func (mapper UpperMapper) Table2Obj(name string) string {
	newstr := make([]rune, 0)
	for i, chr := range name {
		if i > 0 && chr > 'Z' {
			chr += ('z' - 'a')
		}
		newstr = append(newstr, chr)
	}
	return string(newstr)
}

func (mapper UpperMapper) TableName(t string) string {
	return t
}
