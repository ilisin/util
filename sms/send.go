package sms

import (
	"errors"
	"fmt"
)

type Sender interface {
	Send(msg *Message) error
}

// Send sends sms using the given Sender.
func Send(s Sender, msg ...*Message) error {
	for i, m := range msg {
		if err := send(s, m); err != nil {
			return fmt.Errorf("gosms: could not send sms %d: %v", i+1, err)
		}
	}

	return nil
}

func send(s Sender, msg *Message) error {
	if err := s.Send(msg); err != nil {
		return err
	}
	return nil
}

func SendMessage(msg ...*Message) error {
	if len(configs) == 0 || len(httpSmserDefaultKey) == 0{
		return errors.New("no sms config")
	}

	httpSmser := configs[httpSmserDefaultKey]
	c := smsNewClient(httpSmserDefaultKey)
	s := &httpSender{c, httpSmser}
	return Send(s, msg...)
}
