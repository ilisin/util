package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"imooly.net/util/sms"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func GetHandle(c *echo.Context) error {
	logrus.Debug(`ss`)
	err := sms.NewJZSms().SendSms(`15050105241`, `尊敬的用户，您本次的验证码为：122334，有效期为5分钟【魔力网】`)
	logrus.Error(err)
	return err
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	GetHandle(nil)
	e.Get("/", GetHandle)
	e.Run(":8080")
}
