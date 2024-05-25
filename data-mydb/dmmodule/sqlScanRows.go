package dmmodule

import (
	"errors"

	"gorm.io/gorm"
	"request-matcher-openai/data-mydb/mydb"
)

func scanRows(db *gorm.DB, tableSQL string, dbParams []interface{}, tableName string) ([]mydb.DataItem, error) {
	callerFunc := GetParentCallerFunctionName()
	parentCallerFunc := GetGrandCallerFunctionName()
	mylogger.Debugf("scanRows, tableName:%v  caller:%v-%v", tableName, callerFunc, parentCallerFunc)
	rows, err := db.Raw(tableSQL, dbParams...).Rows()
	if err != nil {
		return []mydb.DataItem{}, err
	}
	defer rows.Close()
	if mydb.IsValidDataType(tableName) {
		data := []mydb.DataItem{}
		for rows.Next() {
			item, err := mydb.OperateDBItemDBScanRawRows(tableName, rows, db.ScanRows)
			if err != nil {
				return data, err
			}
			data = append(data, item)
		}
		return data, nil
	}
	return []mydb.DataItem{}, errors.New("unknown typename:" + tableName)
}

func scanRows2(db *gorm.DB, tableSQL string, dbParams []interface{}, tableName string) ([]mydb.DataItem, error) {
	// callerFunc := GetParentCallerFunctionName()
	// parentCallerFunc := GetGrandCallerFunctionName()
	// mylogger.Debugf("scanRows, tableName:%v  caller:%v-%v", tableName, callerFunc, parentCallerFunc)
	rows, err := db.Raw(tableSQL, dbParams...).Rows()
	if err != nil {
		return []mydb.DataItem{}, err
	}
	defer rows.Close()

	data := []mydb.DataItem{}
	for rows.Next() {
		item, err := mydb.MyScanDBRow(db, rows, tableName)
		if err != nil {
			mylogger.Errorf("sql: failed to get row:%v, err:%v", tableName, err)
			continue
		}
		data = append(data, item)
	}
	return data, nil
}
