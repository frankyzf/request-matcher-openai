package myauth

import (
	"errors"
	"fmt"
	"github.com/dpapathanasiou/go-recaptcha"
	"request-matcher-openai/data-mydb/export"
	"request-matcher-openai/gocommon/replyutil"

	"request-matcher-openai/data-mydb/mydb"
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

func (p *MyAuth) GetGroupPrivilege(group string) int {
	return export.GroupPrivilege[group]
}

func (p *MyAuth) VerifyGroup(caller mydb.BaseAccount, requiredGroup string) error {
	if caller.Enable == false {
		return fmt.Errorf("Your account is temporarily disabled by the management. Please visitor the management office for further details.")
	}

	privilege := p.GetGroupPrivilege(caller.AccountType)
	requiredPrivilege := p.GetGroupPrivilege(requiredGroup)
	mylogger.Debugf("user privilege: %v, requiredgroup:%v and privilege:%v", privilege, requiredGroup, requiredPrivilege)
	if privilege < requiredPrivilege {
		err := replyutil.AuthPriviledgeError{Message: "user privilege verification failed"}
		return err
	}
	return nil
}
