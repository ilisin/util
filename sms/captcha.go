package sms

import (
	"fmt"
	"math/rand"
	"time"
)

var rander = rand.New(rand.NewSource(time.Now().Unix()))

func MakeCaptcha() string {
	randv := rander.Int31n(99999)
	return fmt.Sprintf("%05d", randv)
}
