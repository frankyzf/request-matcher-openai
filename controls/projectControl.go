package controls

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"request-matcher-openai/data-mydb/export"
	"request-matcher-openai/data-mydb/mydb"
	"request-matcher-openai/gocommon/replyutil"
)

// @Summary GetProjectList Info
// @Description GetProjectList
// @ID GetProjectList
// @Tags Project
// @Accept  json
// @Produce  json
// @Param from query int false "0"
// @Param size query int false "10"
// @Param search_field query string  false "search field"
// @Param keyword query string  false "keyword"
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Security ApiKeyAuth
// @Router /project/list [get]
func GetProjectList(c *gin.Context) {
	data := []mydb.Project{}
	count := 0
	param, err := export.GetTimeParam(c)
	if err == nil {
		caller, err1 := auth.GetCaller(c.GetString("caller_id"), c.GetString("account_type"), c.GetString("password_update_time"))
		err = err1
		if err == nil {
			err = auth.VerifyGroup(caller, "user")
			if err == nil {
				tableFilter, tableParams := "", []interface{}{}
				data, count, err = manage.GetProjectList(param, tableFilter, tableParams)
			}
		}
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

// @Summary GetProjectItem Info
// @Description GetProjectItem
// @ID GetProjectItem
// @Tags Project
// @Accept  json
// @Produce  json
// @Param     id   path    string     true        "id: 1"
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Security ApiKeyAuth
// @Router /project/item/{id} [get]
func GetProjectItem(c *gin.Context) {
	data := mydb.Project{}
	id := c.Param("id")
	caller, err := auth.GetCaller(c.GetString("caller_id"), c.GetString("account_type"), c.GetString("password_update_time"))
	if err == nil {
		err = auth.VerifyGroup(caller, "user")
		if err == nil {
			data, err = manage.GetOneProject(id)
		}
	}
	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, data)
	}
}

// @Summary CreateProjectItem Info
// @Description CreateProjectItem
// @ID CreateProjectItem
// @Tags Project
// @Accept  json
// @Produce  json
// @Param req body  mydb.ProjectMessage true "req"
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Security ApiKeyAuth
// @Router /project/item [post]
func CreateProjectItem(c *gin.Context) {
	var err error
	data := mydb.Project{}
	req := mydb.ProjectMessage{}
	if err = json.NewDecoder(c.Request.Body).Decode(&req); err == nil {
		caller, err1 := auth.GetCaller(c.GetString("caller_id"), c.GetString("account_type"), c.GetString("password_update_time"))
		err = err1
		if err == nil {
			err = auth.VerifyGroup(caller, "operator")
			if err == nil {
				data, err = manage.CreateOneProject(req)
			}
		}
	}
	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, data)
	}
}

// @Summary UpdateProjectItem Info
// @Description UpdateProjectItem
// @ID UpdateProjectItem
// @Tags Project
// @Accept  json
// @Produce  json
// @Param     id   path    string     true        "id: 1"
// @Param req body  mydb.ProjectMessage  true "req"
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Security ApiKeyAuth
// @Router /project/item/{id} [post]
func UpdateProjectItem(c *gin.Context) {
	var err error
	data := mydb.Project{}
	req := mydb.ProjectMessage{}
	id := c.Param("id")
	if err = json.NewDecoder(c.Request.Body).Decode(&req); err == nil {
		caller, err1 := auth.GetCaller(c.GetString("caller_id"), c.GetString("account_type"), c.GetString("password_update_time"))
		err = err1
		if err == nil {
			err = auth.VerifyGroup(caller, "operator")
			if err == nil {
				data = mydb.ConvertMessageToProject(req)
				data, err = manage.UpdateOneProject(id, data)
				data, err = manage.GetOneProject(id)
			}
		}
	}

	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, data)
	}
}

// @Summary DeleteProjectItem Info
// @Description DeleteProjectItem
// @ID DeleteProjectItem
// @Tags Project
// @Accept  json
// @Produce  json
// @Param     id   path    string     true        "id: 1"
// @Success 200 {object} replyutil.ResBody { "success": true, "errorCode":0, "errorMsg":"", "data":null}
// @Failure 500 {object} replyutil.ResBody { "success": false, "errorCode":200, "errorMsg":"err", "data":null}
// @Security ApiKeyAuth
// @Router /project/item/{id} [delete]
func DeleteProjectItem(c *gin.Context) {
	data := mydb.Project{}
	id := c.Param("id")
	caller, err := auth.GetCaller(c.GetString("caller_id"), c.GetString("account_type"), c.GetString("password_update_time"))
	if err == nil {
		err = auth.VerifyGroup(caller, "operator")
		if err == nil {
			err = manage.DeleteOneProject(id)
		}
	}
	if err != nil {
		replyutil.ResAppErr(c, err)
	} else {
		replyutil.ResOk(c, data)
	}
}
