package dmmodule

import (
	"request-matcher-openai/data-mydb/mydb"
)

func (p *DMManage) UTGetList(tableName string, param mydb.RequestParam, userType string,
	tableFilter string, tableParams []interface{}) ([]mydb.DataItem, int, error) {
	return p.UTGetListByDB(tableName, param, userType, tableFilter, tableParams)
}

func (p *DMManage) UTGetListByDB(tableName string, param mydb.RequestParam, userType string,
	tableFilter string, tableParams []interface{}) ([]mydb.DataItem, int, error) {
	if userType != "" {
		tableFilter += " and user_type_int=? "
		tableParams = append(tableParams, mydb.GetUserTypeInt(userType))
	}
	return p.GetListByDB(tableName, param, tableFilter, tableParams)
}

func (p *DMManage) UTGetCount(tableName string, param mydb.RequestParam, userType string,
	tableFilter string, tableParams []interface{}) (int, error) {
	return p.UTGetCountByDB(tableName, param, userType, tableFilter, tableParams)
}

func (p *DMManage) UTGetCountByDB(tableName string, param mydb.RequestParam, userType string,
	tableFilter string, tableParams []interface{}) (int, error) {
	if userType != "" {
		tableFilter += " and user_type_int=? "
		tableParams = append(tableParams, mydb.GetUserTypeInt(userType))
	}
	return p.GetCountByDB(tableName, param, tableFilter, tableParams)
}

func (p *DMManage) UTSplitGetListByDB(tableName string, param mydb.RequestParam, userType string,
	tableFilter string, tableParams []interface{}) ([]mydb.DataItem, error) {
	if userType != "" {
		tableFilter += " and user_type_int=? "
		tableParams = append(tableParams, mydb.GetUserTypeInt(userType))
	}
	return p.SplitGetListByDB(tableName, param, tableFilter, tableParams)
}

func (p *DMManage) UTSplitGetListByDBWithOrderBy(tableName string, param mydb.RequestParam, userType string,
	tableFilter string, tableParams []interface{}, orderBy string) ([]mydb.DataItem, error) {
	if userType != "" {
		tableFilter += " and user_type_int=? "
		tableParams = append(tableParams, mydb.GetUserTypeInt(userType))
	}
	return p.SplitGetListByDBWithOrderBy(tableName, param, tableFilter, tableParams, orderBy)
}

func (p *DMManage) UTGetListByDBWithOrderBy(tableName string, param mydb.RequestParam,
	userType string, tableFilter string, tableParams []interface{}, orderBy string) ([]mydb.DataItem, int, error) {
	if userType != "" {
		tableFilter += " and user_type_int=? "
		tableParams = append(tableParams, mydb.GetUserTypeInt(userType))
	}
	return p.GetListByDBWithOrderBy(tableName, param, tableFilter, tableParams, orderBy)
}

// == full item
func (p *DMManage) UTGetFullItemList(tableName string, param mydb.RequestParam, userType string,
	tableFilter string, tableParams []interface{},
	joinFilter string, joinParams []interface{}) ([]mydb.DataItem, int, error) {
	return p.UTGetFullItemListByDB(tableName, param, userType, tableFilter, tableParams, joinFilter, joinParams)
}

func (p *DMManage) UTGetFullItemListWithOrderBy(tableName string, param mydb.RequestParam, userType string,
	tableFilter string, tableParams []interface{}, orderBy string,
	joinFilter string, joinParams []interface{}, joinOrderBy string) ([]mydb.DataItem, int, error) {
	return p.UTGetFullItemListByDBWithOrderBy(tableName, param, userType, tableFilter, tableParams, orderBy, joinFilter, joinParams, joinOrderBy)
}

func (p *DMManage) UTGetFullItemListByDB(tableName string, param mydb.RequestParam, userType string,
	tableFilter string, tableParams []interface{},
	joinFilter string, joinParams []interface{}) ([]mydb.DataItem, int, error) {
	if userType != "" {
		tableFilter += " and user_type_int=? "
		tableParams = append(tableParams, mydb.GetUserTypeInt(userType))
	}
	p.mylogger.Infof("UTGetFullItemListByDB, table:%v, tableFilter:%v and tableParam:%v, joinFilter:%v and joinParam:%s",
		tableName, tableFilter, tableParams, joinFilter, joinParams)
	return p.GetFullItemList(tableName, param, tableFilter, tableParams, joinFilter, joinParams)
}

func (p *DMManage) UTGetFullItemListByDBWithOrderBy(tableName string, param mydb.RequestParam, userType string,
	tableFilter string, tableParams []interface{}, orderBy string,
	joinFilter string, joinParams []interface{}, joinOrderBy string) ([]mydb.DataItem, int, error) {
	if userType != "" {
		tableFilter += " and user_type_int=? "
		tableParams = append(tableParams, mydb.GetUserTypeInt(userType))
	}
	p.mylogger.Infof("UTGetFullItemListByDBWithOrderBy, table:%v, userType:%v, tableFilter:%v and tableParam:%v, joinFilter:%v and joinParam:%s, orderBy:%v and joinOrderBy:%v",
		tableName, userType, tableFilter, tableParams, joinFilter, joinParams, orderBy, joinOrderBy)
	data, err := p.SplitGetFullItemListByDBWithOrderBy(tableName, param, userType, tableFilter, tableParams, orderBy, joinFilter, joinParams, joinOrderBy)
	if err != nil {
		p.mylogger.Errorf("failed to getDBJoin:%v, err:%v", tableName, err)
		return data, 0, err
	}
	count, err2 := p.UTGetFullCount(tableName, param, userType, tableFilter, tableParams, joinFilter, joinParams)
	return data, count, err2
}

func (p *DMManage) UTSplitGetFullItemListByDB(tableName string, param mydb.RequestParam, userType string,
	tableFilter string, tableParams []interface{},
	joinFilter string, joinParams []interface{}) ([]mydb.DataItem, error) {
	if userType != "" {
		tableFilter += " and user_type_int=? "
		tableParams = append(tableParams, mydb.GetUserTypeInt(userType))
	}
	return p.SplitGetFullItemListByDB(tableName, param, userType, tableFilter, tableParams, joinFilter, joinParams)
}

func (p *DMManage) UTGetFullCount(tableName string, param mydb.RequestParam, userType string,
	tableFilter string, tableParams []interface{},
	joinFilter string, joinParams []interface{}) (int, error) {
	return p.UTGetFullCountByDB(tableName, param, userType, tableFilter, tableParams, joinFilter, joinParams)
}

func (p *DMManage) UTGetFullCountByDB(tableName string, param mydb.RequestParam, userType string,
	tableFilter string, tableParams []interface{},
	joinFilter string, joinParams []interface{}) (int, error) {
	if userType != "" {
		tableFilter += " and user_type_int=? "
		tableParams = append(tableParams, mydb.GetUserTypeInt(userType))
	}
	return p.GetFullCountByDB(tableName, param, userType, tableFilter, tableParams, joinFilter, joinParams)
}

//==sync

func (p *DMManage) UTGetListByDBForSync(tableName string, param mydb.RequestParam, userType string,
	tableFilter string, tableParams []interface{}, orderBy string) ([]mydb.DataItem, error) {
	if userType != "" {
		tableFilter += " and user_type_int=? "
		tableParams = append(tableParams, mydb.GetUserTypeInt(userType))
	}
	return p.GetListByDBForSync(tableName, param, tableFilter, tableParams, orderBy)
}
