package controls

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"request-matcher-openai/data-common/datamanage"
	"request-matcher-openai/data-common/myauth"
	"request-matcher-openai/data-mydb/export"
	"request-matcher-openai/gocommon/commoncontext"
	"request-matcher-openai/gocommon/util"
)

var mylogger *log.Entry

var myRClient *redis.Client
var auth *myauth.MyAuth
var manage *datamanage.DataManager
var myWebSocketManager export.WebSocketManagerInterface
var myTokenManager *util.TokenManager

func Setup(mauth *myauth.MyAuth, dm *datamanage.DataManager, ws export.WebSocketManagerInterface) error {
	mylogger = commoncontext.SetupLogging("request-matcher-openai", "controls")
	myRClient = commoncontext.GetInstance().RClient
	myTokenManager = commoncontext.GetInstance().MyTokenManager
	auth = mauth
	manage = dm
	myWebSocketManager = ws

	return nil
}

func SetupLoginToken(userID string, accountType string, expireStr string, token string) error {
	if expireStr != "" {
		expire, err := strconv.Atoi(expireStr)
		if err != nil {
			return err
		}
		myTokenManager.SetLoginToken(userID, accountType, int64(expire), token)
		return err
	} else {
		return errors.New("empty expire second")
	}
}

func SendEmail(title, content, email string) error {
	form := url.Values{}
	form.Add("from", commoncontext.GetDefaultString("mail.sender_account", "support<support@mail.homeplus.ai>"))
	form.Add("to", email)
	form.Add("subject", title)
	form.Add("text", content)

	req, err := http.NewRequest("POST",
		commoncontext.GetDefaultString("mail.api_host", "https://api.mailgun.net/v3/mail.homeplus.ai/messages"),
		bytes.NewBufferString(form.Encode()))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return err
	}

	req.SetBasicAuth("api", commoncontext.GetDefaultString("mail.api_key", ""))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return err
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
	return nil
}
