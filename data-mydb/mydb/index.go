package mydb

import (
	"time"
)

type DataItem interface {
	GetID() string
	TableName() string
	GetType() string
	GetUpdateTimestamp() time.Time
}

// three kinds of table, general/user_type_sensitive/status_type_sensitive
func IsUserTypeSensitive(name string) bool {
	return false
}

func IsStatusTypeSensitive(name string) bool {
	return false
}

// 1: resident, 2:visitor, 3: blacklist
func GetUserTypeInt(userType string) int {
	if userType == "resident" || userType == "user" {
		return 1
	} else if userType == "visitor" {
		return 2
	} else if userType == "blacklist" {
		return 3
	}
	return 0
}

func GetNow() time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	return now
}

func GetBeginOfDay(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t = t.In(loc)
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func GetEndOfDay(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t = t.In(loc)
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, t.Location()) //mysql time precision is not high enough

}
