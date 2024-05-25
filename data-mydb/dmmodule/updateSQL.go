package dmmodule

import (
	"errors"

	"gorm.io/gorm"
	"request-matcher-openai/data-mydb/mydb"
)

func createDBOneItem(db *gorm.DB, dbItem mydb.DataItem) (mydb.DataItem, error) {
	callerFunc := GetParentCallerFunctionName()
	parentCallerFunc := GetGrandCallerFunctionName()
	mylogger.Infof("sql:create caller:%v-%v", callerFunc, parentCallerFunc)

	var err error
	db, err = getModelDB(db, dbItem.TableName())
	if err != nil {
		return dbItem, err
	}
	newItem, err2 := mydb.MyOpDBItem(db, dbItem, dbItem.TableName(), "create")
	return newItem, err2
}

func updateDBOneItem(db *gorm.DB, dbItem mydb.DataItem, id string) (mydb.DataItem, error) {
	callerFunc := GetParentCallerFunctionName()
	parentCallerFunc := GetGrandCallerFunctionName()
	mylogger.Infof("sql:update, where_filter:%v and tableParams:%v, caller:%v-%v", "id=?", id, callerFunc, parentCallerFunc)

	newItem, err2 := mydb.MyOpDBItem(db.Where("id=?", id), dbItem, dbItem.TableName(), "update")
	return newItem, err2
}

func updateDBOneItemWithFilter(db *gorm.DB, dbItem mydb.DataItem, whereFilter string, tableParams []interface{}) (mydb.DataItem, error) {
	callerFunc := GetParentCallerFunctionName()
	parentCallerFunc := GetGrandCallerFunctionName()
	mylogger.Infof("sql:update, where_filter:%v and tableParams:%v, caller:%v-%v", whereFilter, tableParams, callerFunc, parentCallerFunc)

	var err error
	db, err = getModelDB(db, dbItem.TableName())
	if err != nil {
		return dbItem, err
	}
	newItem, err2 := mydb.MyOpDBItem(db.Where(whereFilter, tableParams...), dbItem, dbItem.TableName(), "update")
	return newItem, err2
}

func addOrUpdateDBOneItem(db *gorm.DB, dbItem mydb.DataItem, whereFilter string, tableParams []interface{}) (mydb.DataItem, error) {
	callerFunc := GetParentCallerFunctionName()
	parentCallerFunc := GetGrandCallerFunctionName()
	mylogger.Infof("sql:add_or_update, where_filter:%v and tableParams:%v, caller:%v-%v", whereFilter, tableParams, callerFunc, parentCallerFunc)

	var err error
	db, err = getModelDB(db, dbItem.TableName())
	if err != nil {
		return dbItem, err
	}
	newItem, err2 := mydb.MyOpDBItem(db.Where(whereFilter, tableParams...), dbItem, dbItem.TableName(), "first_or_create")
	return newItem, err2
}

func deleteDBItem(db *gorm.DB, tableName string, whereFilter string, tableParams []interface{}) error {
	callerFunc := GetParentCallerFunctionName()
	parentCallerFunc := GetGrandCallerFunctionName()
	mylogger.Infof("sql:delete, where_filter:%v and tableParams:%v,  caller:%v-%v", whereFilter, tableParams, callerFunc, parentCallerFunc)

	var err error
	db, err = getModelDB(db, tableName)
	if err != nil {
		return err
	}
	if rawItem, ok := mydb.GetAllDataTypes()[tableName]; ok {
		item := rawItem
		_, err = mydb.MyOpDBItem(db.Where(whereFilter, tableParams...), item, tableName, "delete")
		return err
	} else {
		return errors.New("unknown typename:" + tableName)
	}

}

// expose it to outside
func SyncDBOneItem(db *gorm.DB, dbItem mydb.DataItem, id string) (mydb.DataItem, error) {
	callerFunc := GetParentCallerFunctionName()
	parentCallerFunc := GetGrandCallerFunctionName()
	mylogger.Infof("sql:sync, where_filter:%v and tableParams:%v, caller:%v-%v", "id=?", id, callerFunc, parentCallerFunc)

	var err error
	db, err = getModelDB(db.Unscoped(), dbItem.TableName())
	if err != nil {
		return dbItem, err
	}
	newItem, err2 := mydb.MyOpDBItem(db.Where("id=?", id), dbItem, dbItem.TableName(), "save")
	return newItem, err2
}

// == update column
func updateDBOneItemFlag(db *gorm.DB, tableName string, id string, fields []string, flags []bool) error {
	callerFunc := GetParentCallerFunctionName()
	parentCallerFunc := GetGrandCallerFunctionName()
	mylogger.Infof("sql:update, where_filter:%v and tableParams:%v, caller:%v-%v", "id=?", id, callerFunc, parentCallerFunc)

	var err error
	if len(fields) != len(flags) {
		return errors.New("failed to match update flag and fields")
	}
	values := map[string]interface{}{}
	for index, field := range fields {
		values[field] = flags[index]
	}
	db, err = getModelDB(db, tableName)
	if err != nil {
		return err
	}
	db = db.Where("id=?", id)
	err = db.Updates(values).Error
	return err
}

func updateDBOneItemFields(db *gorm.DB, tableName string, id string, values map[string]interface{}) error {
	callerFunc := GetParentCallerFunctionName()
	parentCallerFunc := GetGrandCallerFunctionName()
	mylogger.Infof("sql:update, where_filter:%v and tableParams:%v, caller:%v-%v", "id=?", id, callerFunc, parentCallerFunc)

	var err error
	db, err = getModelDB(db, tableName)
	if err != nil {
		return err
	}
	db = db.Where("id=?", id)
	err = db.Updates(values).Error
	return err
}

func updateDBItemsFieldsWithFilter(db *gorm.DB, tableName string, whereFilter string, tableParams []interface{}, values map[string]interface{}) error {
	callerFunc := GetParentCallerFunctionName()
	parentCallerFunc := GetGrandCallerFunctionName()
	mylogger.Debugf("updateDBItemsFlag, whereFilter:%v and tableParams:%v, len param:%v, caller:%v-%v", whereFilter, tableParams, len(tableParams), callerFunc, parentCallerFunc)

	var err error
	db, err = getModelDB(db, tableName)
	if err != nil {
		return err
	}
	sql := db.Where(whereFilter, tableParams...)
	err = sql.Updates(values).Error
	return err
}
