package mongo

import (
	"gogs.xlh/tools/util/logger"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strings"
)

func Insert(db, collection string, docs ...interface{}) (err error) {
	engine, err := NewMongoEngine(db)
	if err != nil {
		logger.Error(err)
	}
	defer engine.Close()
	err = engine.Insert(collection, docs...)
	return
}

func Delete(db, collection string, query interface{}) (err error) {
	engine, err := NewMongoEngine(db)
	if err != nil {
		logger.Error(err)
	}
	defer engine.Close()
	err = engine.Delete(collection, query)
	return
}

func DeleteAll(db, collection string, query interface{}) (err error) {
	engine, err := NewMongoEngine(db)
	if err != nil {
		logger.Error(err)
	}
	defer engine.Close()
	err = engine.DeleteAll(collection, query)
	return
}

func Update(db, collection string, selector, update interface{}) (err error) {
	engine, err := NewMongoEngine(db)
	if err != nil {
		logger.Error(err)
	}
	defer engine.Close()
	err = engine.Update(collection, selector, update)
	return
}

func Upsert(db, collection string, selector, update interface{}) (changeInfo *mgo.ChangeInfo, err error) {
	engine, err := NewMongoEngine(db)
	if err != nil {
		logger.Error(err)
	}
	defer engine.Close()
	return engine.Upsert(collection, selector, update)
}

func UpdateAll(db, collection string, selector, update interface{}) (err error) {
	engine, err := NewMongoEngine(db)
	if err != nil {
		logger.Error(err)
	}
	defer engine.Close()
	err = engine.UpdateAll(collection, selector, update)
	return
}

func Count(db, collection string, query interface{}) (n int, err error) {
	engine, err := NewMongoEngine(db)
	if err != nil {
		logger.Error(err)
	}
	defer engine.Close()
	n, err = engine.Count(collection, query)
	return
}

func GetOne(db, collection string, query, result interface{}) (err error) {
	engine, err := NewMongoEngine(db)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer engine.Close()
	err = engine.GetOne(collection, query, result)
	return
}

func GetPage(db, collection string, skip, limit int, query, result interface{}) (err error) {
	engine, err := NewMongoEngine(db)
	if err != nil {
		logger.Error(err)
	}
	defer engine.Close()
	err = engine.GetPage(collection, skip, limit, query, result)
	return
}

func GetPageSort(db, collection string, skip, limit int, query, result interface{}, sort ...string) (err error) {
	engine, err := NewMongoEngine(db)
	if err != nil {
		logger.Error(err)
	}
	defer engine.Close()
	err = engine.GetPageSort(collection, skip, limit, query, result, sort...)
	return
}

func GetAll(db, collection string, query, result interface{}) (err error) {
	engine, err := NewMongoEngine(db)
	if err != nil {
		logger.Error(err)
	}
	defer engine.Close()
	err = engine.GetAll(collection, query, result)
	return
}

/*
字段名称id
模糊查询 id_regex
In查询 id_in
*/
func GetCondition(db, collection string, p map[string]interface{}, result interface{}) (err error) {
	engine, err := NewMongoEngine(db)
	if err != nil {
		logger.Error(err)
	}
	defer engine.Close()
	selector := bson.M{}
	var fieldValue interface{}
	for k, v := range p {
		if strings.HasSuffix(k, `_regex`) {
			k = k[:len(k)-len(`_regex`)]
			fieldValue = bson.M{"$regex": v}
		} else {
			fieldValue = v
		}
		selector[k] = fieldValue
	}
	err = engine.GetAll(collection, selector, result)
	return
}
