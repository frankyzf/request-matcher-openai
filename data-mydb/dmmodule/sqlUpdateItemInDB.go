package dmmodule

import (
	"errors"

	"gorm.io/gorm"
	"request-matcher-openai/data-mydb/mydb"
)

func updateItemInDB(db *gorm.DB, dbItem mydb.DataItem) (mydb.DataItem, error) {
	typeName := dbItem.GetType()
	callerFunc := GetParentCallerFunctionName()
	parentCallerFunc := GetGrandCallerFunctionName()
	mylogger.Debugf("updateItemInDB, typeName:%v  caller:%v-%v", typeName, callerFunc, parentCallerFunc)

	//db already have model
	if mydb.IsValidDataType(typeName) {
		var err error
		dbItem, err = mydb.OperateDBItemUpdate(dbItem, typeName, db.Updates)
		mylogger.Debugf("firstOrCreateItemInDB,  dbItem:%T, dbItem:%v and err:%v", dbItem, dbItem, err)
		return dbItem, err
	}
	return dbItem, errors.New("unknown typename:" + typeName)

}
