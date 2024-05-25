package dmmodule

import (
	"errors"

	"gorm.io/gorm"
	"request-matcher-openai/data-mydb/mydb"
)

func createItemInDB(db *gorm.DB, dbItem mydb.DataItem) (mydb.DataItem, error) {
	typeName := dbItem.GetType()
	callerFunc := GetParentCallerFunctionName()
	parentCallerFunc := GetGrandCallerFunctionName()
	mylogger.Debugf("createItemInDB, typeName:%v  caller:%v-%v", typeName, callerFunc, parentCallerFunc)

	//db already have model
	if mydb.IsValidDataType(typeName) {
		var err error
		dbItem, err = mydb.OperateDBItemCreate(dbItem, typeName, db.Create)
		return dbItem, err
	}
	return dbItem, errors.New("unknown typename:" + typeName)
}
