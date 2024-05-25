package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"request-matcher-openai/gocommon/commoncontext"
	"request-matcher-openai/gocommon/replyutil"
	"request-matcher-openai/gocommon/util"
)

// make sure the JWTAuth is used after commoncontext is init
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		claims, err := VerifyAuthorization(authorization, secret)
		if err != nil {
			fmt.Printf("JWTAuth err:%v\n", err)
			replyutil.ResAppErr(c, err)
			c.AbortWithStatus(200)
			return
		}

		//password update time
		passwordUpdateTimeUnixMilli, ok4 := claims["password_update_time"].(string)
		if ok4 == false || passwordUpdateTimeUnixMilli == "" {
			err4 := replyutil.AuthExpireError{Message: "Password verification failed, please login again."}
			fmt.Printf("JWTAuth err:%v\n", err4)
			replyutil.ResAppErr(c, err4)
			c.AbortWithStatus(200)
			return
		}

		//login token verify
		err = VerifyLoginToken(claims, authorization)
		if err != nil {
			fmt.Printf("JWTAuth err:%v\n", err)
			replyutil.ResAppErr(c, err)
			c.AbortWithStatus(200)
			return
		}

		//single session verify
		es, _ := claims["expire"].(string) //account use expire as single key
		err = VerifySingleSession(claims, c.GetHeader("device-id"), es)
		if err != nil {
			fmt.Printf("JWTAuth err:%v\n", err)
			replyutil.ResAppErr(c, err)
			c.AbortWithStatus(200)
			return
		}

		fmt.Printf("%v JWTAuth parse company_id:%v caller_id:%v, name:%v, account_type:%v\n",
			time.Now().Format("2006/01/02 - 15:04:05"),
			claims["company_id"], claims["id"], claims["name"], claims["account_type"])
		c.Set("caller_id", claims["id"])
		c.Set("user_type", claims["user_type"])
		c.Set("account_type", claims["account_type"])
		c.Set("password_update_time", claims["password_update_time"])
		c.Set("group", claims["group"])
		c.Set("claims", claims)
	}
}

func VerifyAuthorization(authorization string, secret string) (map[string]interface{}, error) {
	if len(authorization) == 0 {
		return map[string]interface{}{}, replyutil.AuthError{Message: "missing authorization"}
	}

	claims, err := ParseToken(authorization, secret)
	if err != nil {
		return map[string]interface{}{}, replyutil.AuthError{Message: "failed to parse authorization token:" + err.Error()}
	}

	es, ok := claims["expire"].(string)
	unix := int64(0)
	if ok == false || es == "" {
		return map[string]interface{}{}, replyutil.AuthExpireError{Message: "no expire time"}
	}
	unix, err = strconv.ParseInt(es, 10, 64)
	if err != nil || int64(unix) <= time.Now().Unix() {
		return map[string]interface{}{}, replyutil.AuthExpireError{Message: "authorization expired"}
	}

	callerID, ok1 := claims["id"].(string)
	if ok1 == false || callerID == "" {
		return map[string]interface{}{}, replyutil.AuthError{Message: "no id info in the authorization header"}
	}
	claims["caller_id"] = callerID //both id and caller_id are the same

	userType, ok2 := claims["user_type"].(string)
	if ok2 == false || userType == "" {
		claims["user_type"] = "user"
	}

	accountType, ok3 := claims["account_type"].(string)
	if ok3 == false || accountType == "" {
		claims["account_type"] = "user"
	}

	group, ok3 := claims["group"].(string)
	if ok3 == false || group == "" {
		claims["group"] = "user"
	}

	return claims, nil
}

func VerifyLoginToken(claims map[string]interface{}, authorization string) error {
	companyID, _ := claims["company_id"].(string)
	accountType, _ := claims["account_type"].(string)
	callerID, _ := claims["caller_id"].(string)
	es, _ := claims["expire"].(string)
	unix, _ := strconv.ParseInt(es, 10, 64)

	userTokenVerifyEnable := commoncontext.GetDefaultBool(fmt.Sprintf("%v.cognito.login_token.enable_verify", companyID), commoncontext.GetDefaultBool("cognito.login_token.enable_verify", false))
	accountTokenVerifyEnable := commoncontext.GetDefaultBool(fmt.Sprintf("%v.cognito.login_token.enable_account_verify", companyID), commoncontext.GetDefaultBool("cognito.login_token.enable_account_verify", false))
	if (userTokenVerifyEnable && accountType == "user") ||
		(accountTokenVerifyEnable && accountType != "user") { //default turn off token_verify
		if commoncontext.GetInstance().RClient == nil {
			return replyutil.AuthExpireError{Message: "token service busy"}
		}
		key := util.GetLoginTokenKey(accountType, callerID, unix)
		savedToken, err2 := commoncontext.GetInstance().RClient.Get(key).Result()
		if err2 != nil || savedToken != authorization {
			return replyutil.AuthExpireError{Message: "token expired/mismatch"}
		}
	}
	return nil

}

func VerifySingleSession(claims map[string]interface{}, deviceID string, accountExpire string) error {
	companyID, _ := claims["company_id"].(string)
	accountType, _ := claims["account_type"].(string)
	callerID, _ := claims["caller_id"].(string)
	group, _ := claims["group"].(string)

	userSingleSessionVerify := commoncontext.GetDefaultBool(fmt.Sprintf("%v.cognito.single_session.enable_verify", companyID), commoncontext.GetDefaultBool("cognito.single_session.enable_verify", false))
	accountSingleSessionVerify := commoncontext.GetDefaultBool(fmt.Sprintf("%v.cognito.single_session.enable_account_verify", companyID), commoncontext.GetDefaultBool("cognito.single_session.enable_account_verify", false))
	if (userSingleSessionVerify && accountType == "user") ||
		(accountSingleSessionVerify && accountType != "user") { //default turn off single_session
		singleSessionID := ""
		if accountType == "user" {
			singleSessionID = deviceID
		} else {
			singleSessionID = accountExpire // for account, use expireUnix as device-id(webpage does not have device-id)
		}
		if commoncontext.GetInstance().RClient == nil {
			return replyutil.AuthExpireError{Message: "single session service busy"}
		}
		key := util.GetSingleSessionKey(accountType, callerID)
		val, err := commoncontext.GetInstance().RClient.Get(key).Result()
		if err != nil {
			return replyutil.AuthExpireError{Message: "single session id does not exist"}
		}
		fmt.Printf("the save single session id:%v and header single session(device-id) id:%v", val, deviceID)
		if val != singleSessionID {
			if accountType == "account" && group == "company-admin" {
				fmt.Printf("the company-admin account needs to skip checking the single session id matchs var:%v, singleSessionID:%v", val, singleSessionID)
			} else {
				return replyutil.AuthExpireError{Message: "you will be logged out of this session as your account is logged in on another device."}
			}
		} // else pass the check successfully

	}
	return nil
}

func ParseToken(tokenString string, secret string) (map[string]interface{}, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	if err != nil {
		return map[string]interface{}{}, err
	}
	if token == nil {
		return map[string]interface{}{}, errors.New("token parse is nil")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if token.Valid {
			return claims, nil
		} else {
			return map[string]interface{}{}, errors.New("token validation failed")
		}
	} else {
		return map[string]interface{}{}, errors.New("can not convert to token claims")
	}
}
