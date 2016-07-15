package tools

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Coppyer struct {
	senstive  bool //赋值字段大小写是否敏感，true:敏感  false:不敏感
	anonymous bool //匿名成员赋值,匿名成员平铺(递归)，冲突的话则匹配Base.Field
	extends   bool //struct成员赋值(递归)
}

func NewCoppyer(sens bool, anony bool, ext bool) *Coppyer {
	return &Coppyer{
		senstive:  sens,
		anonymous: anony,
		extends:   ext,
	}
}

func (c Coppyer) GenValues(bean interface{}) (map[string]*reflect.Value, error) {
	val := reflect.ValueOf(bean)
	typ := reflect.TypeOf(bean)
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errors.New("传参错误!")
	}
	m := make(map[string]*reflect.Value)
	etyp := typ.Elem()
	eval := reflect.Indirect(val)
	for i := 0; i < eval.NumField(); i++ {
		field := eval.Field(i)
		atype := etyp.Field(i)
		switch field.Kind() {
		case reflect.Struct:
			if atype.Anonymous && c.anonymous {
				for j := 0; j < field.NumField(); j++ {
					iival := field.Field(j)
					iityp := field.Type().Field(j)
					tempName := iityp.Name
					if !c.senstive {
						tempName = strings.ToUpper(tempName)
					}
					if _, ok := m[tempName]; ok {
						return nil, errors.New("字段名冲突")
					}
					m[tempName] = &iival
				}
			} else if atype.Anonymous || c.extends {
				if _, ok := field.Interface().(time.Time); ok { //时间类型不转
					tempName := atype.Name
					if !c.senstive {
						tempName = strings.ToUpper(tempName)
					}
					if _, ok := m[tempName]; ok {
						return nil, errors.New("字段名冲突")
					}
					m[tempName] = &field
					continue
				}

				for j := 0; j < field.NumField(); j++ {
					iival := field.Field(j)
					iityp := field.Type().Field(j)
					tempName := atype.Name + "." + iityp.Name
					if !c.senstive {
						tempName = strings.ToUpper(tempName)
					}
					if _, ok := m[tempName]; ok {
						return nil, errors.New("字段名冲突")
					}
					m[tempName] = &iival
				}
			} else {
				tempName := atype.Name
				if !c.senstive {
					tempName = strings.ToUpper(tempName)
				}
				if _, ok := m[tempName]; ok {
					return nil, errors.New("字段名冲突")
				}
				m[tempName] = &field
			}
		default:
			//tempName := field
			tempName := atype.Name
			if !c.senstive {
				tempName = strings.ToUpper(tempName)
			}
			if _, ok := m[tempName]; ok {
				return nil, errors.New("字段名冲突")
			}
			m[tempName] = &field
		}
	}
	return m, nil
}

func (c Coppyer) Copy(src interface{}, dst interface{}) error {
	msrc, err := c.GenValues(src)
	if err != nil {
		return err
	}
	mdst, err := c.GenValues(dst)
	if err != nil {
		return err
	}
	for k, v := range mdst {
		if sv, ok := msrc[k]; ok {
			if v.Type().Kind() != sv.Type().Kind() {
				return errors.New(fmt.Sprintf("%s 字段类型不匹配", sv))
			}
			v.Set(*sv)
		}
	}
	return nil
}
