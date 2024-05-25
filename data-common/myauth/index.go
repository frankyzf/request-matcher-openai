package myauth

import (
	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"request-matcher-openai/gocommon/commoncontext"
	"request-matcher-openai/gocommon/util"
)

var mylogger *logrus.Entry

type MyAuth struct {
	myDbConn  *gorm.DB
	myRClient *redis.Client
	Secret    string
	bVerify   bool
	bCaptcha  bool
	mylogger  *logrus.Entry
}

func Setup() *MyAuth {

	mylogger = commoncontext.SetupLogging("0", "myauth")

	bVerify := commoncontext.GetDefaultBool("cognito.verify", false)
	secret := commoncontext.GetDefaultString("cognito.secret", "mysecret")
	bCaptcha := commoncontext.GetDefaultBool("captcha.verify", false)
	captchaSecret := commoncontext.GetDefaultString("captcha.secret", "mysecret")

	dbconn := commoncontext.GetInstance().DBConn
	rclient := commoncontext.GetInstance().RClient
	mauth := GetMyAuth(dbconn, rclient, bVerify, secret, bCaptcha, captchaSecret)

	return mauth
}

func GetMyAuth(dbconn *gorm.DB, rclient *redis.Client, verify bool, secret string,
	bCaptcha bool, captchaSecret string) *MyAuth {
	p := &MyAuth{}
	p.myDbConn = dbconn
	p.myRClient = rclient
	p.Secret = secret
	p.bVerify = verify
	p.bCaptcha = bCaptcha
	recaptcha.Init(captchaSecret)
	p.mylogger = commoncontext.SetupLogging("0", "myauth")

	return p
}

func (p MyAuth) VerifySignatureCode(code string) bool {
	signature := commoncontext.GetDefaultString("signature.value", "")
	key := commoncontext.GetDefaultString("signature.key", "")
	signatureCheck := util.VerifySignatureCodeWithSha256(signature, key, code)
	if !signatureCheck {
		p.mylogger.Errorf("error to verify signature code:%v", code)
		return false
	}
	return true
}
