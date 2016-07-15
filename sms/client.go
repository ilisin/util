package sms

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"github.com/Sirupsen/logrus"
)

type httpClient interface {
	Read(req *http.Request) ([]byte, error)
	ParseResult(data []byte) error
}

type SmsClient struct {
	client *http.Client
}

type JZClient struct {
	SmsClient
}

var (
	smsNewClient = func(ty string) httpClient {
		switch ty {
		case HS_JZ:
			return NewJZClient()
		default:
			return NewSmsClient()
		}
	}
)

func NewSmsClient() *SmsClient {
	return &SmsClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func NewJZClient() *JZClient {
	return &JZClient{
		SmsClient: SmsClient{
			client: &http.Client{
				Timeout: 10 * time.Second,
			},
		},
	}
}

func (c *SmsClient) Read(req *http.Request) ([]byte, error) {
	now := time.Now()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(`sms do error[%v], time[%v]`, err.Error(), time.Now().Sub(now)))
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(data) == 0 {
		return nil, errors.New(`sms readall error:` + err.Error())
	}
	return data, nil
}

func (c *SmsClient) ParseResult(data []byte) error {
	return nil
}

func (c *JZClient) ParseResult(data []byte) error {
	smsResult := bytes.NewBuffer(data)
	smsResultInt, err := strconv.Atoi(smsResult.String())
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Debug(smsResultInt)
	if smsResultInt <= 0 {
		return errors.New(`sms not result:` + smsResult.String())
	}
	return nil
}
