package controls

import (
	"encoding/json"
	"errors"
	"request-matcher-openai/gocommon/commoncontext"
	"time"

	"github.com/gin-gonic/gin"

	"request-matcher-openai/data-mydb/mydb"
	"request-matcher-openai/gocommon/model"
	"request-matcher-openai/gocommon/replyutil"
)

// @Summary LoginUser Info
// @Description LoginUser
// @ID LoginUser
// @Tags Login
// @Accept  json
// @Produce  json
// @Param login   body       mydb.UserLoginMessage   true "login"
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Router /app-login [post]
func LoginUser(c *gin.Context) {
	var err error
	data := mydb.UserLoginMessage{}
	item := mydb.UserShort{}
	if err = json.NewDecoder(c.Request.Body).Decode(&data); err == nil {
		uniqID := data.Email
		if data.Email != "" {
			uniqID = data.Email
			item, err = manage.UserLoginHandler(data)
		} else if data.Phone != "" {
			uniqID = data.Phone
			item, err = manage.UserLoginWithPhoneHandler(data)
		}
		if err == nil {
			claims := auth.ProduceUserProfile(item, "user")
			token := ""
			token, err = auth.ProduceToken(claims)
			if err == nil {
				item.Token = token
				item.APIKey = commoncontext.GetDefaultString("dify.api_key", "api_key")
				SetupLoginToken(item.ID, "user", claims["expire"], token)

				myTokenManager.InsertToken(uniqID, model.LoginToken{
					ID:                item.ID,
					Email:             data.Email,
					Phone:             data.Phone,
					Timestamp:         time.Now(),
					Expire:            24 * 3600,
					Token:             token,
					PasswordTimestamp: item.PasswordUpdateTime,
				})
			}
		}
	}
	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, item)
	}
}

// @Summary LogoutUser Info
// @Description LogoutUser
// @ID LogoutUser
// @Tags Login
// @Accept  json
// @Produce  json
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Security ApiKeyAuth
// @Router /logout [post]
func LogoutUser(c *gin.Context) {
	data := map[string]interface{}{}
	caller, err := auth.GetCaller(c.GetString("caller_id"), c.GetString("account_type"), c.GetString("password_update_time"))
	if err == nil {
		if caller.AccountType != "user" {
			err = errors.New("should be user logout interface")
		}
		if err == nil {
			myTokenManager.DeleteToken(caller.Email)
			myTokenManager.DeleteToken(caller.Phone)
		}
	}
	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, data)
	}
}

// @Summary SignupUser Info
// @Description SignupUser
// @ID SignupUser
// @Tags Login
// @Accept  json
// @Produce  json
// @Param signup   body       mydb.UserSignupMessage     true "signup"
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Router /app-signup [post]
func SignupUser(c *gin.Context) {
	var err error
	data := mydb.UserSignupMessage{}
	mydata := mydb.User{}
	item := mydb.UserShort{}
	if err == nil {
		if err = json.NewDecoder(c.Request.Body).Decode(&data); err == nil {
			mydata, err = manage.UserSignupHandler(data)
			item = mydb.ConvertUserToUserShort(mydata)
			if err == nil {
				claims := auth.ProduceUserProfile(item, "user")
				token, err2 := auth.ProduceToken(claims)
				data.Token = token
				data.Password = ""
				err = err2
			}
		}
	}
	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, item)
	}
}

// @Summary LoginAccount Info
// @Description LoginAccount
// @ID LoginAccount
// @Tags Login
// @Accept  json
// @Produce  json
// @Param login   body       mydb.UserLoginMessage   true "login"
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Router /login [post]
func LoginAccount(c *gin.Context) {
	var err error
	data := mydb.UserLoginMessage{}
	item := mydb.UserShort{}
	if err = json.NewDecoder(c.Request.Body).Decode(&data); err == nil {
		item, err = manage.AccountLoginHandler(data, c.ClientIP())
		if err == nil {
			claims := auth.ProduceUserProfile(item, "account")
			token := ""
			token, err = auth.ProduceToken(claims)
			if err == nil {
				item.Token = token
				SetupLoginToken(item.ID, "account", claims["expire"], token)

				myTokenManager.InsertToken(item.Email, model.LoginToken{
					ID:                item.ID,
					Email:             item.Email,
					Phone:             item.Phone,
					Timestamp:         time.Now(),
					Expire:            24 * 3600,
					Token:             token,
					PasswordTimestamp: item.PasswordUpdateTime,
				})
			}
		}
	}
	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, item)
	}
}

// @Summary Signup Info
// @Description Signup
// @ID Signup
// @Tags Login
// @Accept  json
// @Produce  json
// @Param signup   body       mydb.AccountMessage     true "signup"
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Router /signup [post]
func Signup(c *gin.Context) {
	var err error
	data := mydb.AccountMessage{}
	mydata := mydb.Account{}
	item := mydb.UserShort{}
	if err == nil {
		if err = json.NewDecoder(c.Request.Body).Decode(&data); err == nil {
			mydata, err = manage.CreateOneAccount(data)
			item = mydb.ConvertBaseAccountToUserShort(mydb.ConvertAccountToBaseAccount(mydata))
			if err == nil {
				claims := auth.ProduceUserProfile(item, "account")
				token, err2 := auth.ProduceToken(claims)
				data.Token = token
				data.Password = ""
				err = err2
			}
		}
	}
	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, item)
	}
}

// @Summary Logout Info
// @Description Logout
// @ID Logout
// @Tags Login
// @Accept  json
// @Produce  json
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Security ApiKeyAuth
// @Router /logout [post]
func Logout(c *gin.Context) {
	data := map[string]interface{}{}
	caller, err := auth.GetCaller(c.GetString("caller_id"), c.GetString("account_type"), c.GetString("password_update_time"))
	if err == nil {
		if caller.AccountType != "account" {
			err = errors.New("should be user logout interface")
		}
		if err == nil {
			myTokenManager.DeleteToken(caller.Email)
			myTokenManager.DeleteToken(caller.Phone)
		}
	}
	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, data)
	}
}

// @Summary GetMe Info
// @Description GetMe
// @ID GetMe
// @Tags Login
// @Accept  json
// @Produce  json
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Security ApiKeyAuth
// @Router /me [get]
func GetMe(c *gin.Context) {
	data := mydb.UserShort{}
	callerID := c.GetString("caller_id")
	caller, err := auth.GetCaller(callerID, c.GetString("account_type"), c.GetString("password_update_time"))
	mylogger.Infof("auth get me, err:%v and caller:%v", err, caller)
	if err == nil {
		accountType := c.GetString("account_type")
		if accountType == "user" {
			data, err = manage.GetOneFullUser(caller.ID)
			data.APIKey = commoncontext.GetDefaultString("dify.api_key", "api_key")
		} else { //acount
			account, _ := manage.GetOneAccountByDB(caller.ID, false)
			data = mydb.ConvertBaseAccountToUserShort(mydb.ConvertAccountToBaseAccount(account))
			data.APIKey = commoncontext.GetDefaultString("dify.api_key", "api_key")
		}
	}
	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, data)
	}
}

// @Summary GetMyQualifiedProjectList Info
// @Description GetMyQualifiedProjectList
// @ID GetMyQualifiedProjectList
// @Tags Login
// @Accept  json
// @Produce  json
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Security ApiKeyAuth
// @Router /my-qualified-project-list [get]
func GetMyQualifiedProjectList(c *gin.Context) {
	qualifiedProjectList := []string{}
	callerID := c.GetString("caller_id")
	caller, err := auth.GetCaller(callerID, c.GetString("account_type"), c.GetString("password_update_time"))
	mylogger.Infof("auth get me, err:%v and caller:%v", err, caller)
	if err == nil {
		qualifiedProjectList = commoncontext.GetDefaultStringSlice("dify.qualified_project_list", []string{})
	}
	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, map[string]interface{}{
			"program_list": qualifiedProjectList,
		})
	}
}
