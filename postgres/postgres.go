package postgres

import (
	"fmt"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"reflect"
)

type oracleConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	SslMode  string
	Mapper   core.IMapper
	ShowSQL  bool
}

var configs map[string]oracleConfig

func init() {
	configs = make(map[string]oracleConfig)
}

func RegisterPostgres(pname string, host, port, db, user, pwd, mode string, showSQL bool) {
	if _, ok := configs[pname]; ok {
		panic(fmt.Sprintf("Postgres注册重复 [%v]", pname))
	}
	config := oracleConfig{
		Host:     host,
		Port:     port,
		Database: db,
		User:     user,
		Password: pwd,
		SslMode:  mode,
		Mapper:   new(UpperMapper),
		ShowSQL:  showSQL,
	}
	configs[pname] = config
}

func NewPostgresXormInit(pname string) (*xorm.Engine, error) {
	var (
		err            error
		postgresEngine *xorm.Engine
	)

	c, ok := configs[pname]
	if !ok {
		panic(fmt.Sprintf("Oracle数据库未注册 [%v]", pname))
	}

	postgresEngine, err = xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password,
		c.Host, c.Port, c.Database, c.SslMode))
	if err != nil {
		//panic(err)
		return nil, err
	}
	//默认Mapper采用SnakMapper列名为 _ 链接组合
	//postgresEngine.SetMapper(core.NewCacheMapper(c.Mapper))
	postgresEngine.ShowSQL = c.ShowSQL

	//	f, err := os.Create("d:\\sql.log")
	//	if err != nil {
	//		utility.Error(err)
	//	} else {
	//		engine.Logger = xorm.NewSimpleLogger(f)
	//	}

	return postgresEngine, nil
}

func And(oldCondition, condition string, args ...interface{}) string {
	c := len(args)
	if c > 0 {
		for _, v := range args {
			typ := reflect.TypeOf(v)
			val := reflect.ValueOf(v)
			switch typ.Kind() { //多选语句switch
			case reflect.String:
				if len(val.String()) == 0 {
					return oldCondition
				}
			case reflect.Int32, reflect.Int, reflect.Int64:
				if val.Int() == 0 {
					return oldCondition
				}
			default:
				return oldCondition
			}
		}
	}

	if len(oldCondition) > 0 {
		oldCondition += ` AND `
	}
	return oldCondition + condition
}

func Where(condition, sql string) string {
	if len(condition) > 0 {
		sql += ` WHERE ` + condition
	}
	return sql
}

func Limit(sql string, skip, limit int32) string {
	if limit > 0 {
		sql += fmt.Sprintf(` LIMIT %d OFFSET %d`, limit, skip)
	}
	return sql
}
