package oracle

import (
	"fmt"
	"testing"
	"time"
)

type Account struct {
	Id         string `xorm:"pk"`
	Name       string
	Age        int
	Amount     float32
	CreateTime time.Time `xorm:"column(CREATEDATE) created"`
}

func (a Account) TableName() string {
	return "XORM_ACCOUNT"
}

func (a *Account) BeforeInsert() {
	fmt.Println("before insert: %s", a.Name)
}

func (a *Account) AfterInsert() {
	fmt.Println("after insert: %s", a.Name)
}

func TestOracle(t *testing.T) {
	RegisterOracle("testoracle", "192.168.10.82", "1521", "md", "mb_md_user", "mb_md_user", false)
	engine, err := NewOracleXormInit("testoracle")
	if err != nil {
		t.Errorf("初始化xorm错误[%v]", err)
	}
	defer engine.Close()
	n, err := engine.Count(new(Account))
	if err != nil {
		t.Errorf("查询错误[%v]", err)
	}
	t.Log(n)
}
