package myauth

import (
	"strconv"
	"time"

	"request-matcher-openai/data-mydb/mydb"
	"request-matcher-openai/gocommon/commoncontext"
)

func getProfileField(data map[string]interface{}, fieldName string, defaultValue string) string {
	if v, ok := data[fieldName]; ok {
		if ss, ok := v.(string); ok {
			return ss
		}
	}
	return defaultValue
}

func (p *MyAuth) GetProfileFromToken(data map[string]interface{}) mydb.BaseAccount {
	accountType, ok := data["accountType"].(string)
	if !ok {
		accountType = "user"
	}
	item := mydb.BaseAccount{
		ID:          getProfileField(data, "id", ""),
		Name:        getProfileField(data, "name", ""),
		Email:       getProfileField(data, "email", ""),
		AccountType: accountType,
	}
	return item
}

// accountType: user or account
func (p *MyAuth) ProduceUserProfile(item mydb.UserShort, accountType string) map[string]string {
	data := map[string]string{}
	data["id"] = item.ID
	data["name"] = item.Name
	data["email"] = item.Email
	data["phone"] = item.Phone
	data["account_type"] = accountType
	expire := commoncontext.GetDefaultInt("cognito.login_token.user_expire", 24*3600*365)
	data["expire"] = strconv.FormatInt(time.Now().Unix()+int64(expire), 10)
	data["password_update_time"] = strconv.FormatInt(item.PasswordUpdateTime.UnixMilli(), 10)
	return data
}
