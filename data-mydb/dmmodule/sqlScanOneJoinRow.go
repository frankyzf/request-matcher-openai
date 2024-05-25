package dmmodule

import (
	"errors"

	"gorm.io/gorm"
	"request-matcher-openai/data-mydb/mydb"
)

// deprecated
func scanOneJoinRow(db *gorm.DB, joinSQL string, dbParams []interface{}, tableName string) (mydb.DataItem, error) {
	callerFunc := GetParentCallerFunctionName()
	parentCallerFunc := GetGrandCallerFunctionName()
	mylogger.Debugf("scanOneJoinRow, typeName:%v  caller:%v-%v", tableName, callerFunc, parentCallerFunc)

	if mydb.IsValidDataType(tableName) {
		item, err := mydb.OperateDBItemScanOneJoinRow(tableName, db.Raw(joinSQL, dbParams...).Scan)
		return item, err
	}
	return nil, errors.New("unknown typename:" + tableName)
}

func scanOneJoinRow2(db *gorm.DB, joinSQL string, dbParams []interface{}, tableName string) (mydb.DataItem, error) {
	data, err := scanJoinRows2(db, joinSQL, dbParams, tableName)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return data[0], nil
}
