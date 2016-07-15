package sms

type ISms interface {
	SendSms(phone, message string) error
}
