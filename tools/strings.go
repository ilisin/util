package tools

import (
	"fmt"
	"reflect"
	"strings"
)

func Join(a []string, sep string, prefix, suffix string) string {
	if len(a) == 0 {
		return ""
	}
	if len(a) == 1 {
		return prefix + a[0] + suffix
	}
	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(a[i]) + len(prefix) + len(suffix)
	}

	b := make([]byte, n)
	bp := copy(b, prefix+a[0]+suffix)
	for _, s := range a[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], prefix+s+suffix)
	}
	return string(b)
}

// sb：slice结构
// sep : 连接符
// 返回使用连接符链接的字符串
func JoinStr(sb interface{}, sep string) string {
	typ := reflect.TypeOf(sb)
	val := reflect.ValueOf(sb)
	ss := make([]string, 0)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}
	if typ.Kind() != reflect.Slice {
		return ""
	}
	for i := 0; i < val.Len(); i++ {
		ss = append(ss, fmt.Sprintf("%v", val.Index(i).Interface()))
	}
	return strings.Join(ss, sep)
}
