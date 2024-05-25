package main

import (
	"fmt"
	"gorm.io/gorm"
	"request-matcher-openai/data-common/datamanage"
)

func MyInit(dbconn *gorm.DB) error {
	datamanage.InitFolder()
	if dbconn != nil {
		err := datamanage.InitDatabaseTable(dbconn)
		if err != nil {
			fmt.Printf("failed to migrate database:%v", err)
			return err
		}
	}
	return nil
}
