package myauth

import (
	"strconv"

	"request-matcher-openai/data-mydb/mydb"
	"request-matcher-openai/gocommon/replyutil"
)

func (p *MyAuth) GetCaller(userID string, accountType string, cachePasswordUpdateTimeUnixMilli string) (mydb.BaseAccount, error) {
	var err error
	baseAccount := mydb.BaseAccount{}

	if userID == "" {
		return baseAccount, replyutil.AuthExpireError{Message: "empty user id"}
	}

	if accountType == "user" {
		user := mydb.User{}
		err = p.myDbConn.First(&user, "id=?", userID).Error
		if err != nil {
			p.mylogger.Errorf("GetOneUser for user:%v, err:%v", userID, err)
			return baseAccount, replyutil.AuthExpireError{Message: "no such user:" + userID}
		}

		baseAccount = mydb.ConvertUserToBaseAccount(user)
	} else {
		err = replyutil.AuthExpireError{Message: "unknown account type"}
	}

	if err == nil {
		passwordUpdateTimeUnixMilli := strconv.FormatInt(baseAccount.PasswordUpdateTime.UnixMilli(), 10)

		if cachePasswordUpdateTimeUnixMilli != passwordUpdateTimeUnixMilli {
			p.mylogger.Errorf("GetCaller for user:%v, check password different cache time:%v the token password time:%v",
				userID, cachePasswordUpdateTimeUnixMilli, passwordUpdateTimeUnixMilli)

			err = replyutil.AuthReloginError{Message: "The password has been changed. Please login again."}
		}
	}

	return baseAccount, err
}
