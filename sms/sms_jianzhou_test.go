package sms

import (
	"testing"
)

func TestSendSms(t *testing.T) {
	s := NewJZSms()
	err := s.SendSms(`15050105241`, `111【魔力网】`)
	t.Error(err)
}
