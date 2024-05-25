package dmmodule

import (
	"errors"

	"gorm.io/gorm"
	"request-matcher-openai/data-mydb/mydb"
)

// deprecated
func scanJoinRows(db *gorm.DB, joinSQL string, dbParams []interface{}, tableName string) ([]mydb.DataItem, error) {
	callerFunc := GetParentCallerFunctionName()
	parentCallerFunc := GetGrandCallerFunctionName()
	mylogger.Debugf("scanJoinRows, typeName:%v  caller:%v-%v", tableName, callerFunc, parentCallerFunc)
	rows, err := db.Raw(joinSQL, dbParams...).Rows()
	if err != nil {
		return []mydb.DataItem{}, err
	}
	defer rows.Close()

	if mydb.IsValidDataType(tableName) {
		data := []mydb.DataItem{}
		for rows.Next() {
			item, err := mydb.OperateDBItemDBScanJoinRawRows(tableName, rows, db.ScanRows)
			if err != nil {
				return data, err
			}
			data = append(data, item)
		}
		return data, nil
	}
	return []mydb.DataItem{}, errors.New("unknown typename:" + tableName)
}

func scanJoinRows2(db *gorm.DB, joinSQL string, dbParams []interface{}, tableName string) ([]mydb.DataItem, error) {
	// callerFunc := GetParentCallerFunctionName()
	// parentCallerFunc := GetGrandCallerFunctionName()
	// mylogger.Debugf("scanJoinRows2, typeName:%v  caller:%v-%v", tableName, callerFunc, parentCallerFunc)
	rows, err := db.Raw(joinSQL, dbParams...).Rows()
	if err != nil {
		return []mydb.DataItem{}, err
	}
	defer rows.Close()

	data := []mydb.DataItem{}
	for rows.Next() {
		item, err := mydb.MyScanJoinDBRow(db, rows, tableName)
		if err != nil {
			mylogger.Errorf("sql: failed to get join row:%v, err:%v", tableName, err)
			continue
		}
		data = append(data, item)
	}
	return data, nil
}
