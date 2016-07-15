/**
* 数据存储于内存中
 */
package db

import (
	"sync"
)

type MemoryDB struct {
	Datas map[interface{}]interface{}
	lock  sync.RWMutex
}

func (db *MemoryDB) Set(key, value interface{}) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	if db.Datas == nil {
		db.Datas = make(map[interface{}]interface{})
	}
	db.Datas[key] = value
	return nil
}

func (db *MemoryDB) Get(key interface{}) interface{} {
	db.lock.RLock()
	defer db.lock.RUnlock()
	if db.Datas == nil {
		return nil
	}
	if v, ok := db.Datas[key]; ok {
		return v
	} else {
		return nil
	}
}

func (db *MemoryDB) Delete(key interface{}) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	delete(db.Datas, key)
	return nil
}

func (db *MemoryDB) Flush(conditions ...(func(interface{}) bool)) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	if len(conditions) == 0 {
		db.Datas = nil
		db.Datas = make(map[interface{}]interface{})
	} else {
		for k, v := range db.Datas {
			for _, cf := range conditions {
				if cf(v) {
					delete(db.Datas, k)
					break
				}
			}
		}
	}
	return nil
}

func NewMemoryDB() (*MemoryDB, error) {
	return &MemoryDB{Datas: make(map[interface{}]interface{})}, nil
}
