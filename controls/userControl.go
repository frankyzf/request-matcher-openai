package controls

import (
	"encoding/json"
	"errors"
	_ "golang.org/x/image/webp"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"request-matcher-openai/data-mydb/export"
	"request-matcher-openai/data-mydb/mydb"
	"request-matcher-openai/gocommon/replyutil"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary GetUserList Info
// @Description GetUserList
// @ID GetUserList
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Security ApiKeyAuth
// @Router /user/list [get]
func GetUserList(c *gin.Context) {
	data := []mydb.User{}
	count := 0
	param, err := export.GetTimeParam(c)
	if err == nil {
		data, count, err = manage.GetUserList(param, "", []interface{}{})
	}
	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, map[string]interface{}{
			"total": count,
			"list":  data,
		})
	}
}

// @Summary GetUserItem Info
// @Description GetUserItem
// @ID GetUserItem
// @Tags User
// @Accept  json
// @Produce  json
// @Param     id   path    string     true        "id: 1"
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Security ApiKeyAuth
// @Router /user/item/{id} [get]
func GetUserItem(c *gin.Context) {
	id := c.Param("id")
	data, err := manage.GetOneFullUser(id)

	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, data)
	}
}

// @Summary UpdateMyPassword Info
// @Description UpdateMyPassword
// @ID UpdateMyPassword
// @Tags User
// @Accept  json
// @Produce  json
// @Param user   body       mydb.UserUpdatePasswordMessage   true  "user"
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Security ApiKeyAuth
// @Router /me/update-password [post]
func UpdateMyPassword(c *gin.Context) {
	data := mydb.UserUpdatePasswordMessage{}
	item := mydb.UserShort{}
	callerID := c.GetString("caller_id")
	_, err := auth.GetCaller(callerID, c.GetString("account_type"), c.GetString("password_update_time"))
	if err == nil {
		if err = json.NewDecoder(c.Request.Body).Decode(&data); err == nil {
			accountType := c.GetString("account_type")
			if accountType == "user" {
				item, err = manage.UserHandleUpdatePassword(callerID, data)
			} else {
				err = errors.New("only user can update password")
			}
		}
	}
	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, item)
	}
}

// @Summary SelfDeleteWithPhone Info
// @Description SelfDeleteWithPhone
// @ID SelfDeleteWithPhone
// @Tags User
// @Accept  json
// @Produce  json
// @Param user   body       mydb.UserSelfDeleteMessage   true  "user"
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Security ApiKeyAuth
// @Router /me/self-delete-with-phone [post]
func SelfDeleteWithPhone(c *gin.Context) {
	data := map[string]interface{}{}
	req := mydb.UserSelfDeleteMessage{}
	caller, err := auth.GetCaller(c.GetString("caller_id"), c.GetString("account_type"), c.GetString("password_update_time"))
	if err == nil {
		// err = auth.VerifyGroup(caller, "user")
		// if err == nil {
		if err = json.NewDecoder(c.Request.Body).Decode(&req); err == nil {
			err = manage.UserHandleSelfDeleteWithPhone(caller.ID, req)
			if err == nil {
				myTokenManager.DeleteToken(caller.Phone)
			}
		}
		// }
	}
	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, data)
	}
}

// @Summary SelfDeleteWithEmail Info
// @Description SelfDeleteWithEmail
// @ID SelfDeleteWithEmail
// @Tags User
// @Accept  json
// @Produce  json
// @Param user   body       mydb.UserSelfDeleteMessage   true  "user"
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Security ApiKeyAuth
// @Router /me/self-delete-with-email [post]
func SelfDeleteWithEmail(c *gin.Context) {
	data := map[string]interface{}{}
	req := mydb.UserSelfDeleteMessage{}
	caller, err := auth.GetCaller(c.GetString("caller_id"), c.GetString("account_type"), c.GetString("password_update_time"))
	if err == nil {
		// err = auth.VerifyGroup(caller, "user")
		// if err == nil {
		if err = json.NewDecoder(c.Request.Body).Decode(&req); err == nil {
			err = manage.UserHandleSelfDeleteWithEmail(caller.ID, req)
			if err == nil {
				myTokenManager.DeleteToken(caller.Phone)
			}
		}
		// }
	}
	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, data)
	}
}

// @Summary UpdateUserItem Info
// @Description UpdateUserItem
// @ID UpdateUserItem
// @Tags User
// @Accept  json
// @Produce  json
// @Param     id   path    string     true        "id: 1"
// @Param user   body       mydb.UserUpdateMessage   true  "user"
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Security ApiKeyAuth
// @Router /user/item/{id} [post]
func UpdateUserItem(c *gin.Context) {
	var err error
	data := mydb.UserUpdateMessage{}
	mydata := mydb.User{}
	id := c.Param("id")

	if err = json.NewDecoder(c.Request.Body).Decode(&data); err == nil {
		//caller, err1 := auth.GetCaller(c.GetString("caller_id"), c.GetString("account_type"), c.GetString("password_update_time"))
		//err = err1
		if err == nil {
			if err == nil {
				if err == nil {
					accountType := c.GetString("account_type")
					data.UserID = id

					phone := ""
					bChangePhone := false
					if data.Phone != nil {
						phone = *data.Phone
						bChangePhone = true
					}
					data.Phone = nil // phone is not changed because it have independent workflow
					mydata, err = manage.HandleUserUpdateMessage(data)

					mylogger.Infof("the mydata phone is:%v and data phone:%v", mydata.Phone, phone)
					if err == nil && accountType == "account" && data.ExpiredAt == "" { // backend can change it and it's not about changing the expired_at
						user, _ := manage.GetOneUser(id, true)
						if bChangePhone && phone != user.Phone {
							if phone == "" && user.Email == "" {
								err = errors.New("mobile cannot be cleared when email is empty")
							}
							if err == nil {
								err = manage.UpdateUserPhoneWithoutHookByForce(id, phone)
							}
						}
					}
				}
			}
		}
	}

	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, mydata)
	}
}

// @Summary UpdateUserItemEnableFlag Info
// @Description UpdateUserItemEnableFlag
// @ID UpdateUserItemEnableFlag
// @Tags User
// @Accept  json
// @Produce  json
// @Param     id   path    string     true        "id: 1"
// @Param     enable   query    string     true        "true/false"
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Security ApiKeyAuth
// @Router /user/enable-item/{id} [post]
func UpdateUserItemEnableFlag(c *gin.Context) {
	enable, _ := strconv.ParseBool(c.Query("enable"))
	id := c.Param("id")
	//caller, err := auth.GetCaller(c.GetString("caller_id"), c.GetString("account_type"), c.GetString("password_update_time"))
	user, err := manage.GetOneUser(id, true)
	if err == nil {
		err = manage.UpdateOneUserEnableFlag(id, enable)
	}

	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, user)
	}
}

// @Summary DeleteUserItem Info
// @Description DeleteUserItem
// @ID DeleteUserItem
// @Tags User
// @Accept  json
// @Produce  json
// @Param     id   path    string     true        "id: 1"
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Security ApiKeyAuth
// @Router /user/item/{id} [delete]
func DeleteUserItem(c *gin.Context) {
	data := mydb.User{}
	id := c.Param("id")
	//caller, err := auth.GetCaller(c.GetString("caller_id"), c.GetString("account_type"), c.GetString("password_update_time"))
	err := manage.DeleteOneUser(id)

	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, data)
	}
}
