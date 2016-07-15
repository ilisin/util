package db

import (
	"testing"
)

func TestMemory(t *testing.T) {
	db, _ := NewDB("memory")
	db.Set("key1", 15)
	db.Set(15, 30)
	t.Log(db.Get("key1"))
	t.Log(db.Get(15))
	db.Delete(15)
	t.Log(db.Get(15))
	db.Set("key2", 23)
	db.Set("key3", -4)
	db.Set("key4", -1)
	db.Set(200, []int{1, 2, 3})
	t.Logf("%v", db)
	db.Flush(func(v interface{}) bool {
		i, ok := v.(int)
		if ok {
			if i%2 == 0 {
				return false
			}
		}
		return true
	})
	t.Logf("%v", db)
}
