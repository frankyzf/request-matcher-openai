package datamanage

import (
	"fmt"
	"os"

	"gorm.io/gorm"
	"request-matcher-openai/data-mydb/mydb"
)

func InitFolder() {
	os.MkdirAll("log/", 0777)
	os.MkdirAll("config/", 0777)
	os.MkdirAll("picture/", 0777)
	os.MkdirAll("static/", 0777)
	os.MkdirAll("docs/", 0777)
}

func InitDatabaseTable(dbconn *gorm.DB) error {
	fmt.Printf("start to init database table\n")
	data := []interface{}{}
	for _, dt := range mydb.GetAllDataTypes() {
		data = append(data, dt)
	}
	dbconn.AutoMigrate(data...)
	return nil
}
