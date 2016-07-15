package sqlite

import (
	"fmt"
	//"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

type sqliteConfig struct {
	dbPath  string
	ShowSQL bool
}

var configs map[string]sqliteConfig

func init() {
	configs = make(map[string]sqliteConfig)
}

func RegisterSqlite(sname string, dbPath string, showSQL bool) {
	if _, ok := configs[sname]; ok {
		panic(fmt.Sprintf("Sqlite注册重复 [%v]", sname))
	}
	config := sqliteConfig{
		dbPath:  dbPath,
		ShowSQL: showSQL,
	}
	configs[sname] = config
}

func NewSqliteXormInit(pname string) (*xorm.Engine, error) {
	var (
		err          error
		sqliteEngine *xorm.Engine
	)

	c, ok := configs[pname]
	if !ok {
		panic(fmt.Sprintf("sqlite数据库未注册 [%v]", pname))
	}

	sqliteEngine, err = xorm.NewEngine("sqlite3", c.dbPath)
	if err != nil {
		//panic(err)
		return nil, err
	}
	//默认Mapper采用SnakMapper列名为 _ 链接组合
	//sqliteEngine.SetMapper(core.NewCacheMapper(c.Mapper))
	sqliteEngine.ShowSQL = c.ShowSQL

	return sqliteEngine, nil
}
