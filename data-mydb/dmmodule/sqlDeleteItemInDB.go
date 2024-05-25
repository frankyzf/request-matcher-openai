package dmmodule

import (
	"errors"

	"request-matcher-openai/data-mydb/mydb"
	"gorm.io/gorm"
)

func deleteItemInDB(db *gorm.DB, tableName string) error {
	callerFunc := GetParentCallerFunctionName()
	parentCallerFunc := GetGrandCallerFunctionName()
	mylogger.Debugf("deleteItemInDB, tableName:%v  caller:%v-%v", tableName, callerFunc, parentCallerFunc)

	//db already have model
	if dbItem, ok := mydb.GetAllDataTypes()[tableName]; ok {
		_, err := mydb.OperateDBItemDelete(dbItem, tableName, db.Delete)
		return err
	}
	return errors.New("unknown typename:" + tableName)
}
