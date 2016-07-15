package oracle

import (
	"fmt"

	"gogs.xlh/tools/configuration"

	"github.com/Sirupsen/logrus"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-oci8"
	//	_ "github.com/tgulacsi/goracle/godrv"
	"os"
	"reflect"
)

type OracleConfig struct {
	Host        string `conf:"host"`
	Port        string `conf:"port"`
	Sid         string `conf:"sid"`
	User        string `conf:"user"`
	Password    string `conf:"password"`
	Mapper      core.IMapper
	ShowSQL     bool `conf:"showsql"`
	LogLevel    core.LogLevel
	LogLevelStr string `conf:"loglevel"`
}

var configs map[string]OracleConfig

func init() {
	configs = make(map[string]OracleConfig)
	var preConfig = struct {
		Configs map[string]OracleConfig `conf:"util.oracle,omit"`
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

func RegisterOracle(oname string, host, port, sid, user, pwd string, showSQL bool, level ...core.LogLevel) {
	if _, ok := configs[oname]; ok {
		panic(fmt.Sprintf("Oracle注册重复 [%v]", oname))
	}
	config := OracleConfig{
		Host:     host,
		Port:     port,
		Sid:      sid,
		User:     user,
		Password: pwd,
		//Mapper:   new(UpperMapper),
		ShowSQL:  showSQL,
	}
	if len(level) > 0 {
		config.LogLevel = level[0]
	} else {
		config.LogLevel = core.LOG_INFO
	}
	configs[oname] = config
}

func NewOracleXormInit(oname string) (*xorm.Engine, error) {
	var (
		err          error
		oracleEngine *xorm.Engine
	)
	os.Setenv("NLS_LANG", "AMERICAN_AMERICA.AL32UTF8")

	c, ok := configs[oname]
	if !ok {
		panic(fmt.Sprintf("Oracle not regist [%v]", oname))
	}

	//oci8
	oracleEngine, err = xorm.NewEngine("oci8", fmt.Sprintf(`%s/%s@%s:%s/%s`, c.User, c.Password, c.Host, c.Port, c.Sid))
	if err != nil {
		return nil, err
	}
	//默认Mapper采用SnakMapper列名为 _ 链接组合
	//oracleEngine.SetMapper(core.NewCacheMapper(c.Mapper))
	oracleEngine.ShowSQL(c.ShowSQL)
	oracleEngine.Logger().SetLevel(c.LogLevel)

	//	f, err := os.Create("d:\\sql.log")
	//	if err != nil {
	//		utility.Error(err)
	//	} else {
	//		engine.Logger = xorm.NewSimpleLogger(f)
	//	}

	return oracleEngine, nil
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
			case reflect.Bool:
				if val.Bool() {
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
