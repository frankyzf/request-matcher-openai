package dmmodule

import (
	"errors"

	"gorm.io/gorm"
	"request-matcher-openai/data-mydb/mydb"
)

func firstOrCreateItemInDB(db *gorm.DB, dbItem mydb.DataItem) (mydb.DataItem, error) {
	typeName := dbItem.GetType()
	callerFunc := GetParentCallerFunctionName()
	parentCallerFunc := GetGrandCallerFunctionName()
	mylogger.Debugf("firstOrCreateItemInDB, typeName:%v  caller:%v-%v", typeName, callerFunc, parentCallerFunc)

	//db already have model
	if mydb.IsValidDataType(typeName) {
		var err error
		dbItem, err = mydb.OperateDBItemFirstOrCreate(dbItem, typeName, db.Assign(dbItem).FirstOrCreate)
		mylogger.Debugf("firstOrCreateItemInDB,  dbItem:%T, dbItem:%v and err:%v", dbItem, dbItem, err)
		return dbItem, err
	}
	return dbItem, errors.New("unknown typename:" + typeName)
}
