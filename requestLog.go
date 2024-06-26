package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"request-matcher-openai/gocommon/commoncontext"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if commoncontext.GetDefaultBool("requestlog", false) {
			buf, _ := ioutil.ReadAll(c.Request.Body)
			rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
			rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

			fmt.Println(readBody(rdr1)) // Print request body

			c.Request.Body = rdr2
		}
		c.Next()
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
