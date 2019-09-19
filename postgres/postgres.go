package postgres

import (
	"fmt"
	"reflect"

	"github.com/Sirupsen/logrus"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"github.com/ilisin/configuration"
)

type ShowLevel int

type PostgresConfig struct {
	Host        string `conf:"host"`
	Port        string `conf:"port"`
	Database    string `conf:"database"`
	User        string `conf:"user"`
	Password    string `conf:"password"`
	SslMode     string `conf:"sslmode"`
	Mapper      core.IMapper
	ShowSQL     bool `conf:"showsql"`
	LogLevel    core.LogLevel
	LogLevelStr string `conf:"loglevel"`
}

var configs map[string]PostgresConfig

func init() {
	configs = make(map[string]PostgresConfig)
	var preConfig = struct {
		Configs map[string]PostgresConfig `conf:"util.postgres,omit"`
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
		c.LogLevel = strToLogLevel(c.LogLevelStr)
		configs[k] = c
	}
}

func strToLogLevel(str string) core.LogLevel {
	switch str {
	case "OFF", "off":
		return core.LOG_OFF
	case "ERR", "err", "ERROR", "error":
		return core.LOG_ERR
	case "WARN", "warn", "WARNING", "warning":
		return core.LOG_WARNING
	case "debug", "DEBUG":
		return core.LOG_DEBUG
	default:
		return core.LOG_INFO
	}
}

//shows 没有传值则，debug和info都不打印
//shows1 控制debug , shows2控制info
func RegisterPostgres(pname string, host, port, db, user, pwd, mode string, showsql bool, level ...core.LogLevel) {
	if _, ok := configs[pname]; ok {
		panic(fmt.Sprintf("Postgres注册重复 [%v]", pname))
	}
	config := PostgresConfig{
		Host:     host,
		Port:     port,
		Database: db,
		User:     user,
		Password: pwd,
		SslMode:  mode,
		//Mapper:   new(UpperMapper),
		ShowSQL:  showsql,
	}
	if len(level) > 0 {
		config.LogLevel = level[0]
	} else {
		config.LogLevel = core.LOG_INFO
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
		panic(fmt.Sprintf("postgres数据库未注册 [%v]", pname))
	}

	postgresEngine, err = xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password,
		c.Host, c.Port, c.Database, c.SslMode))
	if err != nil {
		//panic(err)
		return nil, err
	}
	//默认Mapper采用SnakMapper列名为 _ 链接组合
	//postgresEngine.SetMapper(core.NewCacheMapper(c.Mapper))
	postgresEngine.ShowSQL(c.ShowSQL)
	postgresEngine.Logger().SetLevel(c.LogLevel)

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

func Or(oldCondition, condition string, args ...interface{}) string {
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
		oldCondition += ` OR `
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
