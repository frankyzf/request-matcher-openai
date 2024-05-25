package replyutil

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResBody is a common struct for the return msg
type ResBody struct {
	Success   bool        `json:"success"`
	ErrorCode int         `json:"errorCode"`
	ErrorMsg  string      `json:"errorMsg"`
	Data      interface{} `json:"data"`
}

// ResOk is the shortcode for return ok msg
func ResOk(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ResBody{true, 0, "ok", data})
}

// ResErr is the shortcode for the return error msg
func ResErr(c *gin.Context, err error) {
	msg := err.Error()
	code := GetErrorCode(err)
	fmt.Printf("reply error code:%v message:%v\n", code, msg)
	c.JSON(http.StatusInternalServerError, ResBody{false, code, msg, nil})
}

func ResWarning(c *gin.Context, data interface{}, err error) { //warning will return data as well as error_code != 0
	warningMsg := err.Error()
	code := GetErrorCode(err)
	fmt.Printf("reply warning code:%v, message:%v\n", code, warningMsg)
	c.JSON(http.StatusOK, ResBody{true, code, warningMsg, data})
}

func ResAppErr(c *gin.Context, err error) {
	msg := err.Error()
	code := GetErrorCode(err)
	fmt.Printf("reply error code:%v, msg:%v\n", code, msg)
	c.JSON(http.StatusOK, ResBody{false, code, msg, nil})
}

// ==
func ResRaw(c *gin.Context, resBody interface{}) {
	c.JSON(200, resBody)
}

// ==
func ResRawByte(c *gin.Context, data []byte) {
	c.Data(200, "application/json", data)
}
