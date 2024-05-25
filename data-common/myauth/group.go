package myauth

import (
	"errors"
	"github.com/dpapathanasiou/go-recaptcha"
)

func (p *MyAuth) VerifyCaptcha(clientIP string, captcha string) error {
	if p.bCaptcha {
		captchaVerify, _ := recaptcha.Confirm(clientIP, captcha)
		if captchaVerify == false {
			return errors.New("error to verify captcha")
		}
	}
	return nil
}
