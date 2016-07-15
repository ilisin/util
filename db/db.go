package db

import "fmt"

type DBInterface interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{} //返回为nil则表示无值
	Delete(key interface{}) error

	/**
	* 数据清除，如果没有传入参数，则删除所有数据，否则清除满足条件的所有数据
	* 例如
	func CleanOld(in interface{}) bool{
		if in == nil {
			return true
		}
		return false
	}
	Flush(CleanOld)
	*/
	Flush(...(func(interface{}) bool)) error
}

//var dbs = make(map[string]DBInterface)
//
//func Register(name string, db DBInterface) {
//	if _, ok := dbs[name]; ok {
//		panic("DB register twice")
//	}
//	dbs[name] = db
//}

func NewDB(name string) (DBInterface, error) {
	if name == "memory" {
		return NewMemoryDB()
	} else {
		return nil, fmt.Errorf("unkown db")
	}
}
