package paramParser

import (
	"errors"
	"log"
	"reflect"
	"strings"
	"time"
)

var (
	errMapNotReg = errors.New("paramer map not register")
	errKeyNotReg = errors.New("paramer key not register")
	errType      = errors.New("param's type error")

	ErrRequired = errors.New("param is required")
)

var _params = make(map[string]ParamMap)

type ParamOption struct {
	Key          string
	Required     bool
	DefaultValue interface{}
}

type ParamMap map[string]ParamOption

func RegisterMap(pname string, options []ParamOption) {
	if _, ok := _params[pname]; !ok {
		_params[pname] = make(ParamMap)
	}
	for _, v := range options {
		if _, ok := _params[pname][v.Key]; !ok {
			_params[pname][v.Key] = v
		} else {
			log.Printf("%v param re register", v.Key)
		}
	}
}

func Parse(pname string, paramKey string, param interface{}) (interface{}, error) {
	if _, ok := _params[pname]; !ok {
		return nil, errMapNotReg
	}
	var (
		option ParamOption
		ok     bool
	)
	if option, ok = _params[pname][paramKey]; !ok {
		return nil, errKeyNotReg
	}

	typ := reflect.TypeOf(param)
	val := reflect.ValueOf(param)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if val.Type().Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, errType
	}

	opTyp := reflect.TypeOf(option.DefaultValue)
	if opTyp.Kind() == reflect.Ptr {
		opTyp = opTyp.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		ft := typ.Field(i).Type
		fv := val.Field(i)
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		if ft.Kind() != reflect.Map {
			continue
		}
		//map的key必须为string
		if ft.Key().Kind() != reflect.String {
			continue
		}
		ele := ft.Elem()
		if ele.Kind() == reflect.Ptr {
			ele = ele.Elem()
		}
		if ele.Kind() == reflect.Slice && opTyp.Kind() == reflect.Slice {
			if ele.Elem().Name() != opTyp.Elem().Name() {
				continue
			}
		}
		if ele.Name() != opTyp.Name() {
			continue
		}
		v := fv.MapIndex(reflect.ValueOf(paramKey))
		//未设置参数,取默认值
		if v.IsValid() == false || v == reflect.Zero(ele) {
			return parseDefault(pname, paramKey)
		} else {
			return v.Interface(), nil
		}
	}

	log.Printf("%v注册的类型不匹配")
	return nil, errType
}

func parseDefault(pname string, paramKey string) (interface{}, error) {
	if _, ok := _params[pname]; !ok {
		return nil, errMapNotReg
	}
	var (
		option ParamOption
		ok     bool
	)
	if option, ok = _params[pname][paramKey]; !ok {
		return nil, errKeyNotReg
	}
	if option.Required {
		return nil, ErrRequired
	}
	return option.DefaultValue, nil
}

//bool 表示  解析空的 int 值
func ParseInt(pname string, paramKey string, param interface{}) (int, bool, error) {
	v, err := Parse(pname, paramKey, param)
	if err != nil {
		return 0, false, err
	}
	if i, ok := v.(int); ok {
		return i, false, nil
	}
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return 0, false, errType
	}
	var (
		vRet int
		bRet bool
		vTag bool
		bTag bool
	)
	for n := 0; n < typ.NumField(); n++ {
		ttyp := typ.Field(n)
		tval := val.Field(n)
		if strings.ToUpper(ttyp.Name) == "VALUE" && ttyp.Type.Kind() == reflect.Int {
			vTag = true
			vRet, _ = tval.Interface().(int)
		}
		if strings.ToUpper(ttyp.Name) == "ISNULL" && ttyp.Type.Kind() == reflect.Bool {
			bTag = true
			bRet, _ = tval.Interface().(bool)
		}
	}
	if vTag && bTag {
		return vRet, bRet, nil
	}
	return 0, false, errType
}

//bool 表示  解析空的 int32 值
func ParseInt32(pname string, paramKey string, param interface{}) (int32, bool, error) {
	v, err := Parse(pname, paramKey, param)
	if err != nil {
		return 0, false, err
	}
	if i, ok := v.(int32); ok {
		return i, false, nil
	}
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return 0, false, errType
	}
	var (
		vRet int32
		bRet bool
		vTag bool
		bTag bool
	)
	for n := 0; n < typ.NumField(); n++ {
		ttyp := typ.Field(n)
		tval := val.Field(n)
		if strings.ToUpper(ttyp.Name) == "VALUE" && ttyp.Type.Kind() == reflect.Int32 {
			vTag = true
			vRet = int32(tval.Int())
		}
		if strings.ToUpper(ttyp.Name) == "ISNULL" && ttyp.Type.Kind() == reflect.Bool {
			bTag = true
			bRet, _ = tval.Interface().(bool)
		}
	}
	if vTag && bTag {
		return vRet, bRet, nil
	}
	return 0, false, errType
}

//bool 表示  解析空的 int32 值
func ParseInt64(pname string, paramKey string, param interface{}) (int64, bool, error) {
	v, err := Parse(pname, paramKey, param)
	if err != nil {
		return 0, false, err
	}
	if i, ok := v.(int64); ok {
		return i, false, nil
	}
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return 0, false, errType
	}
	var (
		vRet int64
		bRet bool
		vTag bool
		bTag bool
	)
	for n := 0; n < typ.NumField(); n++ {
		ttyp := typ.Field(n)
		tval := val.Field(n)
		if strings.ToUpper(ttyp.Name) == "VALUE" && ttyp.Type.Kind() == reflect.Int64 {
			vTag = true
			vRet, _ = tval.Interface().(int64)
		}
		if strings.ToUpper(ttyp.Name) == "ISNULL" && ttyp.Type.Kind() == reflect.Bool {
			bTag = true
			bRet, _ = tval.Interface().(bool)
		}
	}
	if vTag && bTag {
		return vRet, bRet, nil
	}
	return 0, false, errType
}

func ParseFloat64(pname string, paramKey string, param interface{}) (float64, bool, error) {
	v, err := Parse(pname, paramKey, param)
	if err != nil {
		return 0, false, err
	}
	if i, ok := v.(float64); ok {
		return i, false, nil
	}
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return 0, false, errType
	}
	var (
		vRet float64
		bRet bool
		vTag bool
		bTag bool
	)
	for n := 0; n < typ.NumField(); n++ {
		ttyp := typ.Field(n)
		tval := val.Field(n)
		if strings.ToUpper(ttyp.Name) == "VALUE" && ttyp.Type.Kind() == reflect.Float64 {
			vTag = true
			vRet, _ = tval.Interface().(float64)
		}
		if strings.ToUpper(ttyp.Name) == "ISNULL" && ttyp.Type.Kind() == reflect.Bool {
			bTag = true
			bRet, _ = tval.Interface().(bool)
		}
	}
	if vTag && bTag {
		return vRet, bRet, nil
	}
	return 0, false, errType
}

//bool 表示  解析空的 int 值
func ParseString(pname string, paramKey string, param interface{}) (string, bool, error) {
	v, err := Parse(pname, paramKey, param)
	if err != nil {
		return "", false, err
	}
	if s, ok := v.(string); ok {
		return s, false, nil
	}
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	var (
		sRet string
		bRet bool
		sTag bool
		bTag bool
	)
	for n := 0; n < typ.NumField(); n++ {
		ttyp := typ.Field(n)
		tval := val.Field(n)
		if strings.ToUpper(ttyp.Name) == "VALUE" && ttyp.Type.Kind() == reflect.String {
			sTag = true
			sRet, _ = tval.Interface().(string)
		}
		if strings.ToUpper(ttyp.Name) == "ISNULL" && ttyp.Type.Kind() == reflect.Bool {
			bTag = true
			bRet, _ = tval.Interface().(bool)
		}
	}
	if sTag && bTag {
		return sRet, bRet, nil
	}
	return "", false, errType
}

// string list 解析
func ParseStringList(pname string, paramKey string, param interface{}) ([]string, bool, error) {
	v, err := Parse(pname, paramKey, param)
	if err != nil {
		return nil, false, err
	}
	if s, ok := v.([]string); ok {
		return s, false, nil
	}
	return nil, false, errType
}

// int list 解析
func ParseIntList(pname string, paramKey string, param interface{}) ([]int, bool, error) {
	v, err := Parse(pname, paramKey, param)
	if err != nil {
		return nil, false, err
	}
	if is, ok := v.([]int); ok {
		return is, false, nil
	}
	if iis, ok := v.([]int32); ok {
		irs := make([]int, 0)
		for _, i := range iis {
			irs = append(irs, int(i))
		}
		return irs, false, nil
	}
	return nil, false, errType
}

// int list 解析
func ParseBoolList(pname string, paramKey string, param interface{}) ([]bool, bool, error) {
	v, err := Parse(pname, paramKey, param)
	if err != nil {
		return nil, false, err
	}
	if is, ok := v.([]bool); ok {
		return is, false, nil
	}
	return nil, false, errType
}

//bool 表示  解析空的 int 值
func ParseBool(pname string, paramKey string, param interface{}) (bool, bool, error) {
	v, err := Parse(pname, paramKey, param)
	if err != nil {
		return false, false, err
	}
	if b, ok := v.(bool); ok {
		return b, false, nil
	}
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	var (
		bvRet bool
		bRet  bool
		bvTag bool
		bTag  bool
	)
	for n := 0; n < typ.NumField(); n++ {
		ttyp := typ.Field(n)
		tval := val.Field(n)
		if strings.ToUpper(ttyp.Name) == "VALUE" && ttyp.Type.Kind() == reflect.String {
			bvTag = true
			bvRet, _ = tval.Interface().(bool)
		}
		if strings.ToUpper(ttyp.Name) == "ISNULL" && ttyp.Type.Kind() == reflect.Bool {
			bTag = true
			bRet, _ = tval.Interface().(bool)
		}
	}
	if bvTag && bTag {
		return bvRet, bRet, nil
	}
	return false, false, errType
}

//根据int64位的秒数转换
func ParseTime(pname string, paramKey string, param interface{}) (time.Time, bool, error) {
	v, err := Parse(pname, paramKey, param)
	if err != nil {
		return time.Time{}, false, err
	}
	if i, ok := v.(int64); ok {
		return time.Unix(i, 0), false, nil
	}
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return time.Time{}, false, errType
	}
	var (
		vRet int64
		bRet bool
		vTag bool
		bTag bool
	)
	for n := 0; n < typ.NumField(); n++ {
		ttyp := typ.Field(n)
		tval := val.Field(n)
		if strings.ToUpper(ttyp.Name) == "VALUE" && ttyp.Type.Kind() == reflect.Int64 {
			vTag = true
			vRet = tval.Int()
		}
		if strings.ToUpper(ttyp.Name) == "ISNULL" && ttyp.Type.Kind() == reflect.Bool {
			bTag = true
			bRet, _ = tval.Interface().(bool)
		}
	}
	if vTag && bTag {
		return time.Unix(vRet, 0), bRet, nil
	}
	return time.Time{}, false, errType
}

// time list 解析,一般用来解析2个长度的时间数组，用来标示一个时间区间
func ParseTimeList(pname string, paramKey string, param interface{}) ([]time.Time, bool, error) {
	v, err := Parse(pname, paramKey, param)
	if err != nil {
		return nil, false, err
	}
	tl := make([]time.Time, 0)
	if is, ok := v.([]int64); ok {
		for _, iv := range is {
			tl = append(tl, time.Unix(iv, 0))
		}
		return tl, false, nil
	}
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if typ.Kind() != reflect.Slice {
		return tl, false, errType
	}

	for i := 0; i < val.Len(); i++ {
		if val.Field(i).Type().Kind() != reflect.Int64 {
			return tl, false, errType
		}
		tl = append(tl, time.Unix(val.Field(i).Int(), 0))
	}
	return tl, false, nil
}
