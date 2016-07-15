package tools

import (
	"testing"
	"time"
)

type AStruct struct {
	Id       string
	Name     string
	Age      int
	Amount   float32
	BrithDay time.Time
}

type Base struct {
	Say   string
	Count int
}

type CStruct struct {
	CCol  string
	CTime time.Time
}

type BStruct struct {
	Base
	Name     string
	BrithDay time.Time
	CC       CStruct
}

type EStruct struct {
	CCol string
}

type DStruct struct {
	Name     string
	Say      string
	Count    int
	Brithday time.Time
	CC       EStruct
}

func TestGenValues(t *testing.T) {
	return
	//	copy := NewCoppyer(false, false, false)
	//	copy.GenValues(new(AStruct))

	copy2 := NewCoppyer(false, true, true)
	copy2.GenValues(new(BStruct))

	b := &BStruct{}
	b.Base.Count = 12
	b.Base.Say = "hello"
	b.Name = "ilisin"
	b.BrithDay = time.Now()
	b.CC = CStruct{
		CCol: "ccol",
	}
	dst := &DStruct{}
	//	Debug("%v", dst)
	if err := copy2.Copy(b, dst); err != nil {
		t.Error(err)
	}

	//	Debug("%v", dst)
}
