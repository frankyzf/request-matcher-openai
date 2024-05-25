package dmmodule

import (
	"errors"

	"gorm.io/gorm"
	"request-matcher-openai/data-mydb/mydb"
)

// deprecated
func scanOneRow(db *gorm.DB, tableSQL string, dbParams []interface{}, tableName string) (mydb.DataItem, error) {
	// callerFunc := GetParentCallerFunctionName()
	// parentCallerFunc := GetGrandCallerFunctionName()
	// mylogger.Debugf("scanOneRow, typeName:%v  caller:%v-%v", tableName, callerFunc, parentCallerFunc)

	if mydb.IsValidDataType(tableName) {
		item, err := mydb.OperateDBItemScanOneRow(tableName, db.Raw(tableSQL, dbParams...).Scan)
		return item, err
	}
	return nil, errors.New("unknown typename:" + tableName)

}

func scanOneRow2(db *gorm.DB, tableSQL string, dbParams []interface{}, tableName string) (mydb.DataItem, error) {
	data, err := scanRows2(db, tableSQL, dbParams, tableName)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return data[0], nil
}
