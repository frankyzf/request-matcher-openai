package dmmodule

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"request-matcher-openai/data-mydb/mydb"
)

// split get list for single table
func splitGetFilterAndParamForDBSingleTableList(tableName string, param mydb.RequestParam,
	timeFilterField string, tableFilter string, tableParams []interface{}, fieldMap map[string]string) (string, []interface{}) {
	param = fixRequestParam(param)

	searchField, keyword := parseSearchField(param, fieldMap)

	sql := ""
	sql = fmt.Sprintf(` %v where deleted_at is null `, tableName)

	dbParams := []interface{}{}
	if timeFilterField != "" && param.IsTimeFieldSet {
		sql += fmt.Sprintf(" and %v >= ? and %v < ?", timeFilterField, timeFilterField)
		dbParams = append(dbParams, param.Start, param.End)
	}

	sql += tableFilter
	dbParams = append(dbParams, tableParams...)

	if strings.TrimSpace(searchField) != "" && !strings.Contains(searchField, ".") {
		sql += ` and ` + searchField + ` like  ? `
		dbParams = append(dbParams, keyword)
	}
	return sql, dbParams
}

func splitGetListForDBSingleTable(db *gorm.DB, tableName string, param mydb.RequestParam,
	timeFilterField string, tableFilter string, tableParams []interface{}, tableOrder string,
	fieldMap map[string]string) ([]mydb.DataItem, error) {

	sql, dbParams := splitGetFilterAndParamForDBSingleTableList(tableName, param,
		timeFilterField, tableFilter, tableParams, fieldMap)

	sortOrder, _ := parseSortField(tableName, param, fieldMap)
	if sortOrder != "" {
		tableOrder = sortOrder
	}

	if param.From >= 0 {
		sql += tableOrder + ` limit ? offset ?`
		dbParams = append(dbParams, param.Size, param.From)
	}
	tableSQL := `select * from ` + sql

	// fmt.Printf("tableSQL:%v and param:%v", tableSQL, dbParams)
	data, err2 := scanRows2(db, tableSQL, dbParams, tableName)
	return data, err2
}

func splitCountForDBSingleTable(db *gorm.DB, tableName string, param mydb.RequestParam,
	timeFilterField string, tableFilter string, tableParams []interface{}, tableOrder string,
	fieldMap map[string]string) (int, error) {
	sql, dbParams := splitGetFilterAndParamForDBSingleTableList(tableName, param,
		timeFilterField, tableFilter, tableParams, fieldMap)

	cntsql := sql
	cntDbParams := []interface{}{}
	cntDbParams = append(cntDbParams, dbParams...)
	countSQL := `select count(1) as count from ` + cntsql
	count := struct {
		Count int
	}{}
	// fmt.Printf("countSQL:%v and param:%v", countSQL, cntDbParams)
	err := db.Raw(countSQL, cntDbParams...).Scan(&count).Error
	return count.Count, err
}

// == end for split get single table

// begin split join table

func splitGetFilterAndParamForDBJoinList(tableName string, param mydb.RequestParam,
	timeFilterField string, tableFilter string, tableParams []interface{},
	fieldMap map[string]string) (string, []interface{}) {
	param = fixRequestParam(param)

	searchField, keyword := parseSearchField(param, fieldMap)

	sql := fmt.Sprintf(` %v where deleted_at is null `, tableName)
	dbParams := []interface{}{}
	if param.IsTimeFieldSet && timeFilterField != "" {
		sql += fmt.Sprintf(" and %v >= ? and %v < ?", timeFilterField, timeFilterField)
		dbParams = append(dbParams, param.Start, param.End)
	}

	sql += tableFilter
	dbParams = append(dbParams, tableParams...)

	if strings.TrimSpace(searchField) != "" && !strings.Contains(searchField, ".") { //search field does not contain . means it filter on item table, else, filter on join
		sql += ` and ` + searchField + ` like  ? `
		dbParams = append(dbParams, keyword)
	}
	return sql, dbParams
}

func splitCountDBJoinList(db *gorm.DB, tableName string, param mydb.RequestParam,
	timeFilterField string, tableFilter string, tableParams []interface{}, tableOrder string,
	joinSelect string, joinTable string, joinFilter1 string, joinParams []interface{}, joinOrder string,
	fieldMap map[string]string) (int, error) {
	sql, dbParams := splitGetFilterAndParamForDBJoinList(tableName, param, timeFilterField, tableFilter, tableParams, fieldMap)

	param = fixRequestParam(param)
	searchField, keyword := parseSearchField(param, fieldMap)
	count := struct {
		Count int
	}{}

	if strings.TrimSpace(joinFilter1) == "" && !strings.Contains(searchField, ".") {
		cntsql := sql
		cntDbParams := []interface{}{}
		cntDbParams = append(cntDbParams, dbParams...)
		countSQL := `select count(1) as count from ` + cntsql
		err := db.Raw(countSQL, cntDbParams...).Scan(&count).Error //count first
		return count.Count, err
	} else {
		cntsql := sql
		cntDbParams := []interface{}{}
		cntDbParams = append(cntDbParams, dbParams...)
		cntjoinSql := fmt.Sprintf(` ( select * from %v ) as %v `, cntsql, tableName)
		cntjoinSql += joinTable + joinFilter1
		cntDbParams = append(cntDbParams, joinParams...)

		countSQL := `select count(1) as count from ` + cntjoinSql
		if strings.TrimSpace(searchField) != "" && strings.Contains(searchField, ".") {
			countSQL += ` and ` + searchField + ` like  ? `
			cntDbParams = append(cntDbParams, keyword)
		}
		err := db.Raw(countSQL, cntDbParams...).Scan(&count).Error //count first
		return count.Count, err
	}
}

func splitGetListForDBJoin(db *gorm.DB, tableName string, param mydb.RequestParam,
	timeFilterField string, tableFilter string, tableParams []interface{}, tableOrder string,
	joinSelect string, joinTable string, joinFilter1 string, joinParams []interface{}, joinOrder string,
	fieldMap map[string]string) ([]mydb.DataItem, error) {
	sql, dbParams := splitGetFilterAndParamForDBJoinList(tableName, param, timeFilterField, tableFilter, tableParams, fieldMap)

	sortOrder, sortOrderLimitSizeOnSingleTable := parseSortField(tableName, param, fieldMap)
	if sortOrder != "" {
		if sortOrderLimitSizeOnSingleTable {
			tableOrder = sortOrder
			joinOrder = "" // just sort by one
		} else {
			joinOrder = sortOrder
			tableOrder = "" // just sort by one
		}
	}

	param = fixRequestParam(param)
	searchField, keyword := parseSearchField(param, fieldMap)
	bLimitSizeOnSingleTable := true
	if strings.TrimSpace(searchField) != "" && strings.Contains(searchField, ".") ||
		strings.TrimSpace(joinFilter1) != "" {
		bLimitSizeOnSingleTable = false
	}

	if param.From >= 0 {
		if bLimitSizeOnSingleTable {
			sql += tableOrder + ` limit ? offset ?`
			dbParams = append(dbParams, param.Size, param.From)
		} else {
			if sortOrder != "" && sortOrderLimitSizeOnSingleTable {
				sql += tableOrder
			}
		}
	}

	joinSql := fmt.Sprintf(` ( select * from %v ) as %v `, sql, tableName)
	joinSql += joinTable + joinFilter1
	dbParams = append(dbParams, joinParams...)
	if strings.TrimSpace(searchField) != "" && strings.Contains(searchField, ".") {
		joinSql += ` and ` + searchField + ` like  ? `
		dbParams = append(dbParams, keyword)
	}
	if param.From >= 0 && !bLimitSizeOnSingleTable {
		joinSql += joinOrder + ` limit ? offset ?`
		dbParams = append(dbParams, param.Size, param.From)
	} else {
		joinSql += joinOrder
	}
	tableSQL := joinSelect + ` from ` + joinSql
	data, err2 := scanJoinRows2(db, tableSQL, dbParams, tableName)
	return data, err2
}

func splitGetListForDBJoinOnMultiTableAndLimitEndForce(db *gorm.DB, tableName string, param mydb.RequestParam,
	timeFilterField string, tableFilter string, tableParams []interface{}, tableOrder string,
	joinSelect string, joinTable string, joinFilter1 string, joinParams []interface{}, joinOrder string,
	fieldMap map[string]string) ([]mydb.DataItem, error) {
	sql, dbParams := splitGetFilterAndParamForDBJoinList(tableName, param, timeFilterField, tableFilter, tableParams, fieldMap)

	sortOrder, _ := parseSortField(tableName, param, fieldMap)
	if sortOrder != "" {
		joinOrder = sortOrder
	}

	param = fixRequestParam(param)
	searchField, keyword := parseSearchField(param, fieldMap)

	joinSql := fmt.Sprintf(` ( select * from %v ) as %v `, sql, tableName)
	joinSql += joinTable + joinFilter1
	dbParams = append(dbParams, joinParams...)
	if strings.TrimSpace(searchField) != "" && strings.Contains(searchField, ".") {
		joinSql += ` and ` + searchField + ` like  ? `
		dbParams = append(dbParams, keyword)
	}
	if param.From >= 0 {
		joinSql += joinOrder + ` limit ? offset ?`
		dbParams = append(dbParams, param.Size, param.From)
	} else {
		joinSql += joinOrder
	}
	tableSQL := joinSelect + ` from ` + joinSql
	mylogger.Infof("sql:%v, params:%v", tableSQL, tableParams)
	data, err2 := scanJoinRows2(db, tableSQL, dbParams, tableName)
	return data, err2
}
