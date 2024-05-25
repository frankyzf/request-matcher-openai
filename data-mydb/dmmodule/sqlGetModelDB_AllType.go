package dmmodule

import (
	"errors"
	"request-matcher-openai/data-mydb/mydb"

	"gorm.io/gorm"
)

func getModelDB(db *gorm.DB, typeName string) (*gorm.DB, error) {
	if typeName == "user" {
		item := mydb.User{}
		return db.Model(&item), nil
	} else if typeName == "account_user" {
		item := mydb.Account{}
		return db.Model(&item), nil
	} else if typeName == "project" {
		item := mydb.Project{}
		return db.Model(&item), nil
	} else if typeName == "qualification_match" {
		item := mydb.QualificationMatch{}
		return db.Model(&item), nil
	}
	return nil, errors.New("unknown DB typename:" + typeName)
}
