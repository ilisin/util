package sms

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type JZSms struct {
	Url      string
	Account  string
	Password string
}

func NewJZSms() ISms {
	return &JZSms{
		Url:      `http://www.jianzhou.sh.cn/JianzhouSMSWSServer/http/sendBatchMessage`,
		Account:  `sdk_moli`,
		Password: `p9iT9WO3Jw`,
	}
}

func NewJZSmsBy(url, account, password string) ISms {
	return &JZSms{
		Url:      url,
		Account:  account,
		Password: password,
	}
}

/*
发短信
@param : 	phone 手机号 批量用英文;隔开
			message 发送内容
*/
func (s *JZSms) SendSms(phone, message string) error {
	if len(phone+message) == 0 {
		return errors.New(`param error`)
	}

	val := url.Values{}
	val.Set(`account`, s.Account)
	val.Set(`password`, s.Password)
	val.Set(`destmobile`, phone)
	val.Set(`msgText`, message)
	req, _ := http.NewRequest("POST", s.Url, strings.NewReader(val.Encode()))
	req.Header.Add(`Content-Type`, `application/x-www-form-urlencoded;charset=UTF-8`)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	now := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf(`sms do error[%v], time[%v]`, err.Error(), time.Now().Sub(now)))
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(data) == 0 {
		return errors.New(`sms readall error:` + err.Error())
	}
	smsResult := bytes.NewBuffer(data)
	smsResultInt, err := strconv.Atoi(smsResult.String())
	if err != nil {
		return err
	}
	if smsResultInt <= 0 {
		return errors.New(`sms not result:` + smsResult.String())
	}
	return nil
}
