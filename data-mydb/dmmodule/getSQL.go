package dmmodule

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"request-matcher-openai/data-mydb/mydb"
)

//DM_XXXX only care about the database data and cache

func getFieldName(name string, m map[string]string) string {
	if _, ok := m[name]; !ok {
		return name
	}
	return m[name]
}

func fixRequestParam(param mydb.RequestParam) mydb.RequestParam {
	if param.Size == 0 {
		param.Size = 20
	}
	if param.IsTimeFieldSet == false {
		param.Start = time.Unix(0, 0)
		param.End, _ = time.Parse("2006-01-02", "2049-12-31")
		// param.IsTimeFieldSet = true
	}
	return param
}

func parseSearchField(param mydb.RequestParam, fieldMap map[string]string) (string, string) {
	searchField := getFieldName(param.SearchField, fieldMap)
	keyword := param.Keyword
	if keyword == "" {
		searchField = ""
	} else {
		keyword = "%" + keyword + "%"
	}
	return searchField, keyword
}

func parseSortField(tableName string, param mydb.RequestParam, fieldMap map[string]string) (string, bool) {
	sortField := getFieldName(param.SortField, fieldMap)
	sortType := strings.TrimSpace(param.SortType)
	if strings.TrimSpace(sortField) == "" || sortType == "" || (sortType != "asc" && sortType != "desc") {
		return "", true
	}

	if !strings.Contains(sortField, ".") {
		sortField = fmt.Sprintf("%s.%s", tableName, sortField)
	}

	tableOrder := " order by " + sortField + " " + sortType + " "

	return tableOrder, false
}

func getDBSingleTableList(db *gorm.DB, tableName string, param mydb.RequestParam,
	timeFilterField string, tableFilter string, tableParams []interface{}, tableOrder string,
	fieldMap map[string]string) ([]mydb.DataItem, int, error) {
	param = fixRequestParam(param)

	searchField, keyword := parseSearchField(param, fieldMap)

	sortOrder, _ := parseSortField(tableName, param, fieldMap)
	if sortOrder != "" {
		tableOrder = sortOrder
	}

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
	cntsql := sql
	cntDbParams := []interface{}{}
	cntDbParams = append(cntDbParams, dbParams...)
	countSQL := `select count(1) as count from ` + cntsql
	count := struct {
		Count int
	}{}
	// fmt.Printf("countSQL:%v and param:%v", countSQL, cntDbParams)
	db.Raw(countSQL, cntDbParams...).Scan(&count)

	if param.From >= 0 {
		sql += tableOrder + ` limit ? offset ?`
		dbParams = append(dbParams, param.Size, param.From)
	}
	tableSQL := `select * from ` + sql

	// fmt.Printf("tableSQL:%v and param:%v", tableSQL, dbParams)
	data, err2 := scanRows2(db, tableSQL, dbParams, tableName)
	return data, count.Count, err2

}

func getDBSingleTableListForSync(db *gorm.DB, tableName string, param mydb.RequestParam,
	timeFilterField string, tableFilter string, tableParams []interface{}, tableOrder string,
	fieldMap map[string]string) ([]mydb.DataItem, error) {
	param = fixRequestParam(param)

	searchField, keyword := parseSearchField(param, fieldMap)

	sortOrder, _ := parseSortField(tableName, param, fieldMap)
	if sortOrder != "" {
		tableOrder = sortOrder
	}

	sql := fmt.Sprintf(` %v where 1=1 `, tableName)

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

	if param.From >= 0 {
		sql += tableOrder + ` limit ? offset ?`
		dbParams = append(dbParams, param.Size, param.From)
	}
	tableSQL := `select * from ` + sql

	// fmt.Printf("tableSQL:%v and param:%v", tableSQL, dbParams)
	data, err2 := scanRows2(db, tableSQL, dbParams, tableName)
	return data, err2

}

func getDBOneItemFromSingleTable(db *gorm.DB, tableName string,
	tableFilter string, tableParams []interface{}, tableOrder string) (mydb.DataItem, error) {
	sql := ""
	sql = fmt.Sprintf(` %v where deleted_at is null `, tableName)
	dbParams := []interface{}{}

	sql += tableFilter + tableOrder
	dbParams = append(dbParams, tableParams...)

	tableSQL := `select * from ` + sql + " limit 1 offset 0 "
	mylogger.Infof("sql:%v, params:%v", tableSQL, dbParams)
	data, err2 := scanOneRow2(db, tableSQL, dbParams, tableName)
	return data, err2
}

func getDBOneItemFromSingleTableForSync(db *gorm.DB, tableName string,
	tableFilter string, tableParams []interface{}, tableOrder string) (mydb.DataItem, error) {
	sql := ""
	sql = fmt.Sprintf(` %v where 1=1 `, tableName)

	dbParams := []interface{}{}

	sql += tableFilter + tableOrder
	dbParams = append(dbParams, tableParams...)

	tableSQL := `select * from ` + sql + " limit 1 offset 0 "
	mylogger.Infof("sql:%v, params:%v", tableSQL, dbParams)
	data, err2 := scanOneRow2(db, tableSQL, dbParams, tableName)
	return data, err2
}

func getDBJoinList(db *gorm.DB, tableName string, param mydb.RequestParam,
	timeFilterField string, tableFilter string, tableParams []interface{}, tableOrder string,
	joinSelect string, joinTable string, joinFilter1 string, joinParams []interface{}, joinOrder string,
	fieldMap map[string]string) ([]mydb.DataItem, int, error) {
	param = fixRequestParam(param)

	searchField, keyword := parseSearchField(param, fieldMap)
	bLimitSizeOnSingleTable := true
	if strings.TrimSpace(searchField) != "" && strings.Contains(searchField, ".") {
		bLimitSizeOnSingleTable = false
	}

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
	count := struct {
		Count int
	}{}

	if strings.TrimSpace(joinFilter1) == "" && !strings.Contains(searchField, ".") {
		cntsql := sql
		cntDbParams := []interface{}{}
		cntDbParams = append(cntDbParams, dbParams...)
		// cntjoinSql := fmt.Sprintf(` ( select * from %v ) as %v `, cntsql, tableName)
		// cntjoinSql += joinFilter
		// cntDbParams = append(cntDbParams, joinParams...)

		countSQL := `select count(1) as count from ` + cntsql

		// fmt.Printf("joined countSQL:%v and param:%v", countSQL, cntDbParams)
		db.Raw(countSQL, cntDbParams...).Scan(&count) //count first
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

		// fmt.Printf("joined countSQL:%v and param:%v", countSQL, cntDbParams)
		db.Raw(countSQL, cntDbParams...).Scan(&count) //count first
	}

	if param.From >= 0 && bLimitSizeOnSingleTable && strings.TrimSpace(joinFilter1) == "" {
		sql += tableOrder + ` limit ? offset ?`
		dbParams = append(dbParams, param.Size, param.From)
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
	mylogger.Infof("sql:%v params:%v", tableSQL, dbParams)
	// fmt.Printf("joined tableSQL:%v and param:%v", tableSQL, dbParams)
	data, err2 := scanJoinRows2(db, tableSQL, dbParams, tableName)
	return data, count.Count, err2
}

func getDBJoinOneItem(db *gorm.DB, tableName string,
	tableFilter string, tableParams []interface{},
	joinSelect string, joinTable string, joinFilter1 string, joinParams []interface{}) (mydb.DataItem, error) {

	sql := fmt.Sprintf(` %v where deleted_at is null `, tableName)
	dbParams := []interface{}{}

	sql += tableFilter
	dbParams = append(dbParams, tableParams...)

	joinSql := fmt.Sprintf(` ( select * from %v ) as %v `, sql, tableName)
	joinSql += joinTable + joinFilter1
	dbParams = append(dbParams, joinParams...)

	tableSQL := joinSelect + ` from ` + joinSql

	data, err2 := scanOneJoinRow2(db, tableSQL, tableParams, tableName)
	return data, err2
}

func getDBJoinOneItemForSync(db *gorm.DB, tableName string,
	tableFilter string, tableParams []interface{},
	joinSelect string, joinTable string, joinFilter1 string, joinParams []interface{}) (mydb.DataItem, error) {

	sql := fmt.Sprintf(` %v where 1=1 `, tableName)
	dbParams := []interface{}{}

	sql += tableFilter
	dbParams = append(dbParams, tableParams...)

	joinSql := fmt.Sprintf(` ( select * from %v ) as %v `, sql, tableName)
	joinSql += joinTable + joinFilter1
	dbParams = append(dbParams, joinParams...)

	tableSQL := joinSelect + ` from ` + joinSql

	data, err2 := scanOneJoinRow2(db, tableSQL, tableParams, tableName)
	return data, err2
}
