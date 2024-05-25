package mydb

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func OperateDBItemCreate(dbItem DataItem, typeName string, op func(value interface{}) *gorm.DB) (DataItem, error) {
	return opDBItem(dbItem, typeName, op, func(db *gorm.DB) error {
		// if dbItem.GetMyUserTypeInt() > 0 {
		// 	db.UpdateColumn("user_type_int", dbItem.GetMyUserTypeInt())
		// }
		// if dbItem.GetMyDirectionInt() > 0 {
		// 	db.UpdateColumn("direction_int", dbItem.GetMyDirectionInt())
		// }
		return nil
	})
}

func OperateDBItemUpdate(dbItem DataItem, typeName string, op func(value interface{}) *gorm.DB) (DataItem, error) {
	return opDBItem(dbItem, typeName, func(value interface{}) *gorm.DB {
		return op(value)
	}, func(db *gorm.DB) error {
		// if dbItem.GetMyUserTypeInt() > 0 {
		// 	db.UpdateColumn("user_type_int", dbItem.GetMyUserTypeInt())
		// }
		// if dbItem.GetMyDirectionInt() > 0 {
		// 	db.UpdateColumn("direction_int", dbItem.GetMyDirectionInt())
		// }
		return nil
	})
}

func OperateDBItemDelete(dbItem DataItem, typeName string, op func(value interface{}, where ...interface{}) *gorm.DB) (DataItem, error) {
	return opDBItem(dbItem, typeName, func(value interface{}) *gorm.DB {
		return op(value)
	}, func(db *gorm.DB) error {
		err2 := db.Unscoped().UpdateColumn("updated_at", time.Now()).Error
		return err2
	})
}

func OperateDBItemFirstOrCreate(dbItem DataItem, typeName string, op func(value interface{}, where ...interface{}) *gorm.DB) (DataItem, error) {
	return opDBItem(dbItem, typeName, func(value interface{}) *gorm.DB {
		return op(value)
	}, func(db *gorm.DB) error {
		// if dbItem.GetMyUserTypeInt() > 0 {
		// 	db.UpdateColumn("user_type_int", dbItem.GetMyUserTypeInt())
		// }
		// if dbItem.GetMyDirectionInt() > 0 {
		// 	db.UpdateColumn("direction_int", dbItem.GetMyDirectionInt())
		// }
		return nil
	})
}

func OperateDBItemSave(dbItem DataItem, typeName string,
	op func(value interface{}) *gorm.DB, afterHook func(*gorm.DB) error) (DataItem, error) {
	return opDBItem(dbItem, typeName, op, afterHook)
}

func opDBItem(dbItem DataItem, typeName string, op func(value interface{}) *gorm.DB,
	afterHook func(*gorm.DB) error) (DataItem, error) {
	if typeName == "user" {
		item, ok := dbItem.(User)
		if !ok {
			return nil, errors.New("failed to convert dbItem for " + typeName)
		}
		db := op(&item)
		if afterHook != nil {
			afterHook(db)
		}
		return item, db.Error
	} else if typeName == "account_user" {
		item, ok := dbItem.(Account)
		if !ok {
			return nil, errors.New("failed to convert dbItem for " + typeName)
		}
		db := op(&item)
		if afterHook != nil {
			afterHook(db)
		}
		return item, db.Error
	} else if typeName == "project" {
		item, ok := dbItem.(Project)
		if !ok {
			return nil, errors.New("failed to convert dbItem for " + typeName)
		}
		db := op(&item)
		if afterHook != nil {
			afterHook(db)
		}
		return item, db.Error
	} else if typeName == "qualification_match" {
		item, ok := dbItem.(QualificationMatch)
		if !ok {
			return nil, errors.New("failed to convert dbItem for " + typeName)
		}
		db := op(&item)
		if afterHook != nil {
			afterHook(db)
		}
		return item, db.Error
	}
	fmt.Printf("unknown typename:%v\n", typeName)
	return nil, errors.New("unknown typename:" + typeName)
}

func OperateDBItemScanOneRow(typeName string, op func(value interface{}) *gorm.DB) (DataItem, error) {
	return scanOneDBItem(typeName, func(value interface{}) error {
		op(value)
		return nil
	})
}

func OperateDBItemDBScanRawRows(typeName string, rows *sql.Rows, op func(rows *sql.Rows, result interface{}) error) (DataItem, error) {
	return scanOneDBItem(typeName, func(value interface{}) error {
		return op(rows, value)
	})
}

func OperateDBItemScanOneJoinRow(typeName string, op func(value interface{}) *gorm.DB) (DataItem, error) {
	return scanOneFullDBItem(typeName, func(value interface{}) error {
		op(value)
		return nil
	})
}

func OperateDBItemDBScanJoinRawRows(typeName string, rows *sql.Rows, op func(rows *sql.Rows, result interface{}) error) (DataItem, error) {
	return scanOneFullDBItem(typeName, func(value interface{}) error {
		return op(rows, value)
	})
}

func scanOneDBItem(typeName string, op func(value interface{}) error) (DataItem, error) {
	if typeName == "user" {
		item := User{}
		err := op(&item)
		if err != nil {
			return item, err
		}
		if item.GetID() == "" {
			return item, gorm.ErrRecordNotFound
		}
		return item, nil
	} else if typeName == "account_user" {
		item := Account{}
		err := op(&item)
		if err != nil {
			return item, err
		}
		if item.GetID() == "" {
			return item, gorm.ErrRecordNotFound
		}
		return item, nil
	} else if typeName == "project" {
		item := Project{}
		err := op(&item)
		if err != nil {
			return item, err
		}
		if item.GetID() == "" {
			return item, gorm.ErrRecordNotFound
		}
		return item, nil
	} else if typeName == "qualification_match" {
		item := QualificationMatch{}
		err := op(&item)
		if err != nil {
			return item, err
		}
		if item.GetID() == "" {
			return item, gorm.ErrRecordNotFound
		}
		return item, nil
	}
	fmt.Printf("unknown typename:%v\n", typeName)
	return nil, errors.New("unknown typename:" + typeName)

}

func scanOneFullDBItem(typeName string, op func(value interface{}) error) (DataItem, error) {
	if typeName == "user" {
		item := UserShort{}
		err := op(&item)
		if err != nil {
			return item, err
		}
		if item.GetID() == "" {
			return item, gorm.ErrRecordNotFound
		}
		return item, nil
	} else if typeName == "account_user" {
		item := AccountShort{}
		err := op(&item)
		if err != nil {
			return item, err
		}
		if item.GetID() == "" {
			return item, gorm.ErrRecordNotFound
		}
		return item, nil
	} else if typeName == "project" {
		item := ProjectShort{}
		err := op(&item)
		if err != nil {
			return item, err
		}
		if item.GetID() == "" {
			return item, gorm.ErrRecordNotFound
		}
		return item, nil
	} else if typeName == "qualification_match" {
		item := QualificationMatchShort{}
		err := op(&item)
		if err != nil {
			return item, err
		}
		if item.GetID() == "" {
			return item, gorm.ErrRecordNotFound
		}
		return item, nil
	}
	fmt.Printf("unknown typename:%v\n", typeName)
	return nil, errors.New("unknown typename:" + typeName)
}

func GetOwnerFieldName(data DataItem) string {
	return ""
}
