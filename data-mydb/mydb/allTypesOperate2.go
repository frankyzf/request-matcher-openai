package mydb

import (
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func MyScanDBRow(db *gorm.DB, rows *sql.Rows, typeName string) (DataItem, error) {
	if typeName == "user" {
		item := User{}
		db.ScanRows(rows, &item)
		return item, nil
	} else if typeName == "account_user" {
		item := Account{}
		db.ScanRows(rows, &item)
		return item, nil
	}
	fmt.Printf("unknown typename:%v\n", typeName)
	return nil, errors.New("unknown typename:" + typeName)
}

func MyScanJoinDBRow(db *gorm.DB, rows *sql.Rows, typeName string) (DataItem, error) {
	if typeName == "user" {
		item := UserShort{}
		db.ScanRows(rows, &item)
		return item, nil
	} else if typeName == "account_user" {
		item := AccountShort{}
		db.ScanRows(rows, &item)
		return item, nil
	}
	fmt.Printf("unknown typename:%v\n", typeName)
	return nil, errors.New("unknown typename:" + typeName)
}

func MyOpDBItem(db *gorm.DB, dbItem DataItem, typeName string, operate string) (DataItem, error) {
	if typeName == "user" {
		item, ok := dbItem.(User)
		if !ok {
			return nil, errors.New("failed to convert dbItem for " + typeName)
		}
		if operate == "create" {
			err := db.Create(&item).Error
			return item, err
		} else if operate == "first_or_create" {
			err := db.Assign(item).FirstOrCreate(&item).Error
			return item, err
		} else if operate == "update" {
			err := db.Updates(&item).Error
			return item, err
		} else if operate == "save" {
			updatedAt := item.GetUpdateTimestamp()
			err := db.Save(&item).UpdateColumn("updated_at", updatedAt).Error
			return item, err
		} else if operate == "delete" {
			db.UpdateColumn("updated_at", time.Now())
			err := db.Delete(&item).Error
			return item, err
		} else {
			return item, errors.New("unknown operate:" + operate)
		}
	} else if typeName == "account_user" {
		item, ok := dbItem.(Account)
		if !ok {
			return nil, errors.New("failed to convert dbItem for " + typeName)
		}
		if operate == "create" {
			err := db.Create(&item).Error
			return item, err
		} else if operate == "first_or_create" {
			err := db.Assign(item).FirstOrCreate(&item).Error
			return item, err
		} else if operate == "update" {
			err := db.Updates(&item).Error
			return item, err
		} else if operate == "save" {
			updatedAt := item.GetUpdateTimestamp()
			err := db.Save(&item).UpdateColumn("updated_at", updatedAt).Error
			return item, err
		} else if operate == "delete" {
			db.UpdateColumn("updated_at", time.Now())
			err := db.Delete(&item).Error
			return item, err
		} else {
			return item, errors.New("unknown operate:" + operate)
		}
	}
	fmt.Printf("unknown typename:%v\n", typeName)
	return nil, errors.New("unknown typename:" + typeName)
}
