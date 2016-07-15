package mongo

import (
	"testing"
	"time"
)

type TestData struct {
	Name       string
	Age        int
	CreateTime time.Time
}

func TestOracle(t *testing.T) {
	RegisterMongo("testmongo", "192.168.10.49:27017", "log", "root", "123456")
	engine, err := NewMongoEngine("testmongo")
	if err != nil {
		t.Errorf("初始化xorm错误[%v]", err)
	}
	defer engine.Close()
	err = engine.Insert("test", TestData{"hel", 12, time.Now()})
	if err != nil {
		t.Errorf("插入数据失败 [%v]", err)
	}
}
