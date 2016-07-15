package tools

import (
	"testing"
)

func TestJoinStr(t *testing.T) {
	is := []int{12, 32, 44}
	t.Logf("%v", JoinStr(is, ","))
}
