package sms

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"gogs.xlh/tools/configuration"
	"net/http"
	"net/url"
	"strings"
)

const (
	HS_VALUES_PHONE   string = "PHONE"
	HS_VALUES_CAPTION string = "CAPTION"
)

const (
	HS_JZ = "JZ"
)

type HttpSmser struct {
	Url       string            `conf:"url"`
	Model     string            `conf:"model"`  //GET\POST
	Values    map[string]string `conf:"values"` //values: PHONE,CAPTION则动态替换
	Headers   map[string]string `conf:"headers"`
	Signature string            `conf:"signature"`
	Open bool `conf:"open"`
}

type httpSender struct {
	httpClient
	h *HttpSmser
}

var configs map[string]*HttpSmser
var httpSmserDefaultKey string

func init() {
	configs = make(map[string]*HttpSmser)
	templates = make(map[string]*TemplateItem)
	var preConfig = struct {
		Configs map[string]*HttpSmser `conf:"util.sms.api,omit"`
		TemplateItems map[string]*TemplateItem `conf:"util.sms.template,omit"`
	}{}
	err := configuration.Var(&preConfig)
	if err != nil {
		logrus.Fatal(err)
	}
	if preConfig.Configs == nil {
		return
	}
	for k, c := range preConfig.Configs {
		if len(httpSmserDefaultKey) == 0 && c.Open {
			httpSmserDefaultKey = k
		}

		if _, ok := configs[k]; ok {
			panic(fmt.Sprintf("HttpSmser注册重复 [%v]", k))
		}
		configs[k] = c
	}
	logrus.Info(httpSmserDefaultKey, configs[httpSmserDefaultKey])

	for k, v := range preConfig.TemplateItems {
		if _, ok := templates[k]; ok {
			panic(fmt.Sprintf("Template注册重复 [%v]", k))
		}
		templates[k] = v
		logrus.Info(k, *v)
	}

}

func AddHttpSmser(key string, c *HttpSmser) {
	if _, ok := configs[key]; !ok {
		configs[key] = c
	}
}

func (s *HttpSmser) getUrlValues(msg *Message) url.Values {
	val := url.Values{}
	for k, v := range s.Values {
		if v == HS_VALUES_PHONE {
			v = msg.getPhone()
		} else if v == HS_VALUES_CAPTION {
			v = msg.getCaption() + s.Signature
		}
		val.Set(k, v)
	}
	return val
}

func (s *HttpSmser) getRequest(msg *Message) (*http.Request, error) {
	val := s.getUrlValues(msg)
	req, err := http.NewRequest(s.Model, s.Url, strings.NewReader(val.Encode()))
	if err != nil {
		return nil, err
	}
	for k, v := range s.Headers {
		req.Header.Add(k, v)
	}
	return req, nil
}

func (s *httpSender) Send(msg *Message) error {
	if err := msg.check(); err != nil {
		return err
	}
	req, err := s.h.getRequest(msg)
	if err != nil {
		return err
	}
	data, err := s.Read(req)
	if err != nil {
		return err
	}
	return s.ParseResult(data)
}
