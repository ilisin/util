package sms

import (
	"errors"
	"strings"
	"reflect"
	"strconv"
	"fmt"
	"math/rand"
	"time"
)

type Message struct {
	Phone   string
	Caption string
	TemplateKey string //模板编码（KEY）
	TemplateMap map[string]interface{} //code-value 指定code的值
	//SecurityCode string //验证码 为空则自动生成
}

type TemplateItem struct {
	Title  string `conf:"title"`
	Caption string `conf:"caption"`
}
var templates map[string]*TemplateItem
func (m *Message) check() error {
	if len(m.getPhone()) == 0 || len(m.getCaption()) == 0 {
		return errors.New("param err")
	}
	return nil
}

func (m *Message) getPhone() string {
	return m.Phone
}

func (m *Message) getCaption() string {
	if len(m.Caption) == 0 && len(m.TemplateKey) > 0 {
		m.Caption = m.getTemplateValue()
	}
	return m.Caption
}

func (m *Message) CreateSecurityCode() string {
	var rander = rand.New(rand.NewSource(time.Now().Unix()))
	return fmt.Sprintf("%05d", rander.Int31n(99999))
}

func (m *Message) getTemplateKey() string {
	return m.TemplateKey
}

func (m *Message) getTemplateValue() string {
	var caption string
	if v, ok := templates[m.TemplateKey]; ok {
		caption = v.Caption
		for k, value := range m.TemplateMap {
			typ := reflect.TypeOf(value)
			val := reflect.ValueOf(value)
			switch typ.Kind() {
			case reflect.String:
				caption = strings.Replace(caption, k, val.String(), -1)
			case reflect.Int32, reflect.Int, reflect.Int64:
				caption = strings.Replace(caption, k, strconv.FormatInt(val.Int(), 10), -1)
			case reflect.Float32, reflect.Float64:
				caption = strings.Replace(caption, k, fmt.Sprintf("%.2f", val.Float()), -1)
			}
		}
	}
	return caption
}

func (m *Message) AddTemplateMap(key string, value interface{}) {
	if _, ok := m.TemplateMap[key]; !ok {
		m.TemplateMap[key] = value
	}
}