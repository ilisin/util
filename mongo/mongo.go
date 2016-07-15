package mongo

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"gogs.xlh/tools/configuration"

	"labix.org/v2/mgo"
)

var configs map[string]MongoConfig

type MongoConfig struct {
	Hosts    []string `conf:"host"`
	Database string   `conf:"database"`
	UserName string   `conf:"user"`
	Password string   `conf:"password"`
}

type MongoCall func(*mgo.Collection) error

type MongoEngine struct {
	database        string
	mongoDBDialInfo *mgo.DialInfo
	mongoSession    *mgo.Session
}

func init() {
	configs = make(map[string]MongoConfig)
	var preConfig = struct {
		Configs map[string]MongoConfig `conf:"util.mongo,omit"`
	}{}
	err := configuration.Var(&preConfig)
	if err != nil {
		logrus.Fatal(err)
	}
	if preConfig.Configs == nil {
		return
	}
	for k, c := range preConfig.Configs {
		if _, ok := configs[k]; ok {
			panic(fmt.Sprintf("Postgres注册重复 [%v]", k))
		}
		configs[k] = c
	}
}

func RegisterMongo(mname string, host, db, user, pwd string) {
	if _, ok := configs[mname]; ok {
		panic(fmt.Sprintf("mongodb重复注册 [%v]", mname))
	}
	configs[mname] = MongoConfig{
		Hosts:    []string{host},
		Database: db,
		UserName: user,
		Password: pwd}
}

func NewMongoEngine(mname string) (engine *MongoEngine, err error) {
	c, ok := configs[mname]
	if !ok {
		panic(fmt.Sprintf("数据库未注册 [%v]", mname))
	}

	// Create the database object
	engine = &MongoEngine{
		mongoDBDialInfo: &mgo.DialInfo{
			Addrs:    c.Hosts,
			Timeout:  60 * time.Second,
			Database: c.Database,
			Username: c.UserName,
			Password: c.Password,
		},
		database: c.Database,
	}

	// Establish the master session
	engine.mongoSession, err = mgo.DialWithInfo(engine.mongoDBDialInfo)
	if err != nil {
		return engine, err
	}

	engine.mongoSession.SetMode(mgo.Monotonic, true) //mgo.Strong

	// Have the session check for errors
	// http://godoc.org/labix.org/v2/mgo#Session.SetSafe
	engine.mongoSession.SetSafe(&mgo.Safe{})
	return engine, nil
}

func (engine *MongoEngine) Close() {
	engine.mongoSession.Close()
}

// Execute the MongoDB literal function
func (engine *MongoEngine) Execute(collectionName string, mongoCall MongoCall) (err error) {

	// Capture the specified collection
	collection := engine.mongoSession.DB(engine.database).C(collectionName)

	// Execute the mongo call
	err = mongoCall(collection)
	if err != nil {
		return err
	}
	return err
}

// Execute the MongoDB literal function
func (engine *MongoEngine) DBAction(collectionName string, mongoCall MongoCall) (err error) {
	return engine.Execute(collectionName, mongoCall)
}

func (engine *MongoEngine) Insert(collection string, docs ...interface{}) (err error) {
	err = engine.DBAction(collection, func(collection *mgo.Collection) error {
		return collection.Insert(docs...)
	})
	return
}

func (engine *MongoEngine) Delete(collection string, query interface{}) (err error) {
	err = engine.DBAction(collection, func(collection *mgo.Collection) error {
		return collection.Remove(query)
	})
	return
}

func (engine *MongoEngine) DeleteAll(collection string, query interface{}) (err error) {
	err = engine.DBAction(collection, func(collection *mgo.Collection) error {
		_, err = collection.RemoveAll(query)
		return err
	})
	return
}

func (engine *MongoEngine) Update(collection string, selector, update interface{}) (err error) {
	err = engine.DBAction(collection, func(collection *mgo.Collection) error {
		return collection.Update(selector, update)
	})
	return
}

func (engine *MongoEngine) Upsert(collection string, selector, update interface{}) (changeInfo *mgo.ChangeInfo, err error) {
	err = engine.DBAction(collection, func(collection *mgo.Collection) error {
		changeInfo, err = collection.Upsert(selector, update)
		return err
	})
	return
}

func (engine *MongoEngine) UpdateAll(collection string, selector, update interface{}) (err error) {
	err = engine.DBAction(collection, func(collection *mgo.Collection) error {
		_, err = collection.UpdateAll(selector, update)
		return err
	})
	return
}

func (engine *MongoEngine) Count(collection string, query interface{}) (n int, err error) {
	err = engine.DBAction(collection, func(collection *mgo.Collection) error {
		n, err = collection.Find(query).Count()
		return err
	})
	return
}

func (engine *MongoEngine) GetOne(collection string, query, result interface{}) (err error) {
	err = engine.DBAction(collection, func(collection *mgo.Collection) error {
		return collection.Find(query).One(result)
	})

	//	if err == mgo.ErrNotFound {
	//		err = nil
	//	}
	return
}

func (engine *MongoEngine) GetPage(collection string, skip, limit int, query, result interface{}) (err error) {
	err = engine.DBAction(collection, func(collection *mgo.Collection) error {
		return collection.Find(query).Skip(skip).Limit(limit).All(result)
	})
	return
}

func (engine *MongoEngine) GetPageSort(collection string, skip, limit int, query, result interface{}, sort ...string) (err error) {
	err = engine.DBAction(collection, func(collection *mgo.Collection) error {
		return collection.Find(query).Sort(sort...).Skip(skip).Limit(limit).All(result)
	})
	return
}

func (engine *MongoEngine) GetAll(collection string, query, result interface{}) (err error) {
	err = engine.DBAction(collection, func(collection *mgo.Collection) error {
		return collection.Find(query).All(result)
	})
	return
}
