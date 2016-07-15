package postgres

import (
	"fmt"
	"reflect"
	"time"
)

type Query struct {
	Sql       string
	Condition string
}

func NewQuery() *Query {
	return new(Query)
}

func (q *Query) SetSql(sql string, args ...interface{}) {
	if len(args) > 0 {
		q.Sql = fmt.Sprintf(sql, args...)
	} else {
		q.Sql = sql
	}
}

func (q *Query) AddSql(sql string, args ...interface{}) {
	if len(args) > 0 {
		q.Sql += ` ` + fmt.Sprintf(sql, args...)
	} else {
		q.Sql += ` ` + sql
	}
}

func (q *Query) IsConditionExist() bool {
	return len(q.Condition) > 0
}

/*
	判断参数是否有效
*/
func (q *Query) IsValid(args ...interface{}) bool {
	c := len(args)
	if c > 0 {
		for _, v := range args {
			typ := reflect.TypeOf(v)
			val := reflect.ValueOf(v)
			switch typ.Kind() { //多选语句switch
			case reflect.String:
				if len(val.String()) == 0 {
					return false
				}
			case reflect.Int32, reflect.Int, reflect.Int64:
				if val.Int() == 0 {
					return false
				}
			case reflect.Bool:
				if val.Bool() {
					return false
				}
			default:
				return false
			}
		}
	}
	return true
}

func (q *Query) ParseTimeString(t int64) string {
	return time.Unix(t, 0).Format(`2006-01-02 15:04:05`)
}

func (q *Query) And(condition string, args ...interface{}) {
	if !q.IsValid(args...) {
		return
	}
	if q.IsConditionExist() {
		q.Condition += ` AND `
	}
	q.Condition += condition
	return
}

func (q *Query) Or(condition string, args ...interface{}) {
	if !q.IsValid(args...) {
		return
	}

	if q.IsConditionExist() {
		q.Condition += ` OR `
	}
	q.Condition += condition
	return
}

/*
	组合SQL语句
*/
func (q *Query) MakeSql() string {
	if q.IsConditionExist() {
		q.Sql += ` WHERE ` + q.Condition
	}
	return q.Sql
}

func (q *Query) GetCondition() string {
	return q.Condition
}

func (q *Query) Limit(sql string, skip, limit int32) {
	if limit > 0 {
		q.Sql += fmt.Sprintf(` LIMIT %d OFFSET %d`, limit, skip)
	}
	return
}
