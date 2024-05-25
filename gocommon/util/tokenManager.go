package util

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"

	"request-matcher-openai/gocommon/model"
)

type TokenManager struct {
	RClient  *redis.Client
	Duration time.Duration
}

func GetTokenManager(rclient *redis.Client, durationMinute int) *TokenManager {
	p := &TokenManager{
		RClient:  rclient,
		Duration: time.Duration(durationMinute) * time.Minute,
	}
	return p
}

func GetLoginTokenKey(accountType string, callerID string, expireUnix int64) string {
	key := fmt.Sprintf("login_token:%v:%v:%v", accountType, callerID, expireUnix)
	return key
}

func (p *TokenManager) SetLoginToken(userID string, accountType string, expireUnix int64, token string) error {
	key := GetLoginTokenKey(accountType, userID, expireUnix)
	duration := expireUnix - time.Now().Unix()
	err := p.RClient.Set(key, token, time.Duration(duration)*time.Second).Err()
	fmt.Printf("set login token:%v with duration:%v and token:%v, err:%v\n", key, duration, token, err)
	return err
}

func GetSingleSessionKey(accountType string, userID string) string {
	key := fmt.Sprintf("client_device_id:%v:%v", accountType, userID)
	return key
}

func (p *TokenManager) SetSingleSessionDeviceID(userID string, accountType string, expireUnix int64, deviceID string) error {
	key := GetSingleSessionKey(accountType, userID)
	duration := expireUnix - time.Now().Unix()
	err := p.RClient.Set(key, deviceID, time.Duration(duration)*time.Second).Err()
	fmt.Printf("set single session:%v with duration:%v and device_id:%v, err:%v\n", key, duration, deviceID, err)
	return err
}

// uniqID is either phone or email
func (p *TokenManager) InsertToken(uniqID string, login model.LoginToken) error {
	data, err := json.Marshal(login)
	if err != nil {
		return err
	}
	err = p.RClient.Set(uniqID, data, p.Duration).Err()
	return err
}

func (p *TokenManager) GetToken(uniqID string) (model.LoginToken, error) {
	token := model.LoginToken{}

	val, err := p.RClient.Get(uniqID).Result()

	if err != nil {
		return token, err
	}

	err = json.Unmarshal([]byte(val), &token)
	if err != nil {
		return token, err
	}

	return token, nil
}

func (p *TokenManager) DeleteToken(uniqID string) error {
	err := p.RClient.Del(uniqID).Err()
	if err != nil {
		return err
	}
	return nil
}
