package openaimodule

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"request-matcher-openai/data-mydb/mydb"
	"request-matcher-openai/gocommon/logger"
)

type OpenaiAdapter struct {
	accessTokenURL  string
	liveDataBaseURL string
	mylogger        *logrus.Entry
}

func GetOpenaiAdapter(accessTokenURL, liveDataBaseURL string) *OpenaiAdapter {
	p := &OpenaiAdapter{
		accessTokenURL:  accessTokenURL,
		liveDataBaseURL: liveDataBaseURL,
		mylogger:        logger.SetupLogging("info", "openai_module", "openai_adapter"),
	}
	return p
}

func (p *OpenaiAdapter) doSend(myurl string, method string, payload interface{}, headers map[string]string, contentType string) ([]byte, error) {
	httpClient := &http.Client{Timeout: time.Minute}
	buf := new(bytes.Buffer)
	if payload != nil {
		if contentType == "application/x-www-form-urlencoded" || contentType == "multipart/form-data" {
			values := url.Values{}
			payloadMap, _ := payload.(map[string]string)
			for k, v := range payloadMap {
				values.Add(k, v)
			}
			body := strings.NewReader(values.Encode())
			buf.ReadFrom(body)
		} else {
			payloadBuf, _ := json.Marshal(payload)
			buf = bytes.NewBuffer(payloadBuf)
		}
	}
	req, err := http.NewRequest(method, myurl, buf)
	if err != nil {
		p.mylogger.Errorf("error to get http request:%v", err)
		return []byte{}, err
	}
	if headers != nil && len(headers) > 0 {
		for headerKey, headerVal := range headers {
			req.Header[headerKey] = []string{headerVal}
		}
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	command, _ := mydb.PrintCurlCommand(req)
	p.mylogger.Infof("send command to openai adapter:%v", command)
	resp, err2 := httpClient.Do(req)
	if err2 != nil {
		p.mylogger.Errorf("error in send request:%v", err2)
		return []byte{}, err2
	}
	mydb.PrintResponse(resp)
	buf1, _ := io.ReadAll(resp.Body)
	code := resp.StatusCode
	if code > 299 || code < 200 {
		p.mylogger.Errorf("failed to get openai response, code:%v and err:%v", code, string(buf1))
		return []byte{}, errors.New("failed to get openai response")
	}

	return buf1, nil
}
