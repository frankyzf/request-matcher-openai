package dmmodule

import (
	"errors"

	"gorm.io/gorm"
	"request-matcher-openai/data-mydb/mydb"
)

func saveItemInDB(db *gorm.DB, dbItem mydb.DataItem) (mydb.DataItem, error) {
	typeName := dbItem.GetType()
	callerFunc := GetParentCallerFunctionName()
	parentCallerFunc := GetGrandCallerFunctionName()
	mylogger.Debugf("saveItemInDB, typeName:%v  caller:%v-%v", typeName, callerFunc, parentCallerFunc)
	updateTime := dbItem.GetUpdateTimestamp()

	//db already have model
	if mydb.IsValidDataType(typeName) {
		var err error
		dbItem, err = mydb.OperateDBItemSave(dbItem, typeName, db.Save,
			func(db *gorm.DB) error {
				err2 := db.UpdateColumn("updated_at", updateTime).Error
				return err2
			})
		return dbItem, err
	}
	return dbItem, errors.New("unknown typename:" + typeName)
}
