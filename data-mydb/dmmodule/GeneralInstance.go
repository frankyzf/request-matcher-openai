package dmmodule

import (
	"errors"
	"fmt"
	"strings"

	"request-matcher-openai/data-mydb/mydb"
)

func (p *DMManage) GetList(tableName string, param mydb.RequestParam, tableFilter string,
	tableParams []interface{}) ([]mydb.DataItem, int, error) {

	return p.GetListByDB(tableName, param, tableFilter, tableParams)
}

func (p *DMManage) GetListByDB(tableName string, param mydb.RequestParam,
	tableFilter string, tableParams []interface{}) ([]mydb.DataItem, int, error) {
	data, err := p.SplitGetListByDB(tableName, param, tableFilter, tableParams)
	if err != nil {
		p.mylogger.Errorf("failed to get:%v, err:%v ", tableName, err)
		return nil, 0, err
	}
	count, err := p.GetCount(tableName, param, tableFilter, tableParams)
	return data, count, err
}

func (p *DMManage) GetCount(tableName string, param mydb.RequestParam,
	tableFilter string, tableParams []interface{}) (int, error) {
	return p.GetCountByDB(tableName, param, tableFilter, tableParams)
}

func (p *DMManage) GetCountByDB(tableName string, param mydb.RequestParam,
	tableFilter string, tableParams []interface{}) (int, error) {
	count, err := splitCountForDBSingleTable(p.myDbConn, tableName, param,
		p.GetTimeField(tableName), tableFilter, tableParams,
		p.GetOrderBy(tableName), p.GetFieldMap(tableName))
	if err != nil {
		p.mylogger.Errorf("failed to get:%v, err:%v ", tableName, err)
		return count, err
	}
	return count, nil
}

func (p *DMManage) SplitGetListByDB(tableName string, param mydb.RequestParam,
	tableFilter string, tableParams []interface{}) ([]mydb.DataItem, error) {
	orderBy := p.GetOrderBy(tableName)
	return p.SplitGetListByDBWithOrderBy(tableName, param, tableFilter, tableParams, orderBy)
}

func (p *DMManage) SplitGetListByDBWithOrderBy(tableName string, param mydb.RequestParam,
	tableFilter string, tableParams []interface{}, orderBy string) ([]mydb.DataItem, error) {
	data, err := splitGetListForDBSingleTable(p.myDbConn, tableName, param,
		p.GetTimeField(tableName), tableFilter, tableParams,
		orderBy, p.GetFieldMap(tableName))
	if err != nil {
		p.mylogger.Errorf("failed to get:%v, err:%v ", tableName, err)
		return []mydb.DataItem{}, err
	}
	return data, nil
}

func (p *DMManage) GetListByDBWithOrderBy(tableName string, param mydb.RequestParam,
	tableFilter string, tableParams []interface{}, orderBy string) ([]mydb.DataItem, int, error) {
	data, err := splitGetListForDBSingleTable(p.myDbConn, tableName, param,
		p.GetTimeField(tableName), tableFilter, tableParams, orderBy, p.GetFieldMap(tableName))
	if err != nil {
		p.mylogger.Errorf("failed to get:%v, err:%v ", tableName, err)
		return []mydb.DataItem{}, 0, err
	}
	count, err := p.GetCount(tableName, param, tableFilter, tableParams)
	return data, count, err
}

// === full item
func (p *DMManage) GetFullItemList(tableName string, param mydb.RequestParam,
	tableFilter string, tableParams []interface{},
	joinFilter string, joinParams []interface{}) ([]mydb.DataItem, int, error) {
	return p.GetFullItemListByDB(tableName, param, tableFilter, tableParams, joinFilter, joinParams)
}

func (p *DMManage) GetFullItemListByDB(tableName string, param mydb.RequestParam,
	tableFilter string, tableParams []interface{},
	joinFilter string, joinParams []interface{}) ([]mydb.DataItem, int, error) {
	data, err := p.SplitGetFullItemListByDB(tableName, param, "", tableFilter, tableParams, joinFilter, joinParams)
	// p.mylogger.Infof("GetFullItemListByDB, data:%T, err:%v, data:%v", data, err, data)
	if err != nil {
		p.mylogger.Errorf("failed to getDBJoin:%v, err:%v", tableName, err)
		return data, 0, err
	}
	count, err2 := p.GetFullCount(tableName, param, tableFilter, tableParams, joinFilter, joinParams)
	// p.mylogger.Infof("GetFullItemListByDB, count:%v, err:%v", count, err2)
	return data, count, err2
}

func (p *DMManage) SplitGetFullItemListByDB(tableName string, param mydb.RequestParam, userType string,
	tableFilter string, tableParams []interface{},
	joinFilter string, joinParams []interface{}) ([]mydb.DataItem, error) {
	orderBy := p.GetOrderBy(tableName)
	joinOrderBy := p.GetFullItemOrderBy(tableName)
	return p.SplitGetFullItemListByDBWithOrderBy(tableName, param, userType, tableFilter, tableParams, orderBy, joinFilter, joinParams, joinOrderBy)
}

func (p *DMManage) SplitGetFullItemListByDBOnMultiTableAndLimitEndForce(tableName string, param mydb.RequestParam,
	tableFilter string, tableParams []interface{},
	joinFilter string, joinParams []interface{}) ([]mydb.DataItem, error) {
	orderBy := p.GetOrderBy(tableName)
	joinOrderBy := p.GetFullItemOrderBy(tableName)
	return p.SplitGetFullItemListByDBWithOrderByOnMultiTableAndLimitEndForce(tableName, param, tableFilter, tableParams, orderBy, joinFilter, joinParams, joinOrderBy)
}

func (p *DMManage) SplitGetFullItemListByDBWithOrderBy(tableName string, param mydb.RequestParam,
	userType string, tableFilter string, tableParams []interface{}, orderBy string,
	joinFilter string, joinParams []interface{}, joinOrderBy string) ([]mydb.DataItem, error) {
	if orderBy == "" {
		orderBy = p.GetOrderBy(tableName)
	}
	if joinOrderBy == "" {
		joinOrderBy = p.GetFullItemOrderBy(tableName)
	}
	data, err := splitGetListForDBJoin(p.myDbConn, tableName, param,
		p.GetTimeField(tableName), tableFilter, tableParams, orderBy,
		p.GetJoinSelect(tableName), p.GetTableJoin(tableName, userType), joinFilter, joinParams,
		joinOrderBy, p.GetFieldMap(tableName))
	if err != nil {
		p.mylogger.Errorf("failed to getDBJoin:%v, err:%v", tableName, err)
		return data, err
	}
	return data, nil
}

func (p *DMManage) SplitGetFullItemListByDBWithOrderByOnMultiTableAndLimitEndForce(tableName string, param mydb.RequestParam,
	tableFilter string, tableParams []interface{}, orderBy string,
	joinFilter string, joinParams []interface{}, joinOrderBy string) ([]mydb.DataItem, error) {
	if orderBy == "" {
		orderBy = p.GetOrderBy(tableName)
	}
	if joinOrderBy == "" {
		joinOrderBy = p.GetFullItemOrderBy(tableName)
	}
	data, err := splitGetListForDBJoinOnMultiTableAndLimitEndForce(p.myDbConn, tableName, param,
		p.GetTimeField(tableName), tableFilter, tableParams, orderBy,
		p.GetJoinSelect(tableName), p.GetTableJoin(tableName, ""), joinFilter, joinParams,
		joinOrderBy, p.GetFieldMap(tableName))
	if err != nil {
		p.mylogger.Errorf("failed to getDBJoin:%v, err:%v", tableName, err)
		return data, err
	}
	return data, nil
}

func (p *DMManage) GetFullCount(tableName string, param mydb.RequestParam,
	tableFilter string, tableParams []interface{},
	joinFilter string, joinParams []interface{}) (int, error) {
	return p.GetFullCountByDB(tableName, param, "", tableFilter, tableParams, joinFilter, joinParams)
}

func (p *DMManage) GetFullCountByDB(tableName string, param mydb.RequestParam, userType string,
	tableFilter string, tableParams []interface{},
	joinFilter string, joinParams []interface{}) (int, error) {
	count, err := splitCountDBJoinList(p.myDbConn, tableName, param,
		p.GetTimeField(tableName), tableFilter, tableParams, p.GetOrderBy(tableName),
		p.GetJoinSelect(tableName), p.GetTableJoin(tableName, userType), joinFilter, joinParams,
		p.GetFullItemOrderBy(tableName), p.GetFieldMap(tableName))
	if err != nil {
		p.mylogger.Errorf("failed to getDBJoin:%v, err:%v", tableName, err)
		return count, err
	}
	return count, nil
}

func (p *DMManage) GetFullItemListByDBWithOrderBy(tableName string, param mydb.RequestParam,
	tableFilter string, tableParams []interface{}, orderBy string,
	joinFilter string, joinParams []interface{}, joinOrderBy string) ([]mydb.DataItem, int, error) {
	if orderBy == "" {
		orderBy = p.GetOrderBy(tableName)
	}
	if joinOrderBy == "" {
		joinOrderBy = p.GetFullItemOrderBy(tableName)
	}
	data, err := splitGetListForDBJoin(p.myDbConn, tableName, param,
		p.GetTimeField(tableName), tableFilter, tableParams, orderBy,
		p.GetJoinSelect(tableName), p.GetTableJoin(tableName, ""), joinFilter, joinParams,
		joinOrderBy, p.GetFieldMap(tableName))
	if err != nil {
		p.mylogger.Errorf("failed to getDBJoin:%v, err:%v", tableName, err)
		return data, 0, err
	}
	count, err2 := p.GetFullCount(tableName, param, tableFilter, tableParams, joinFilter, joinParams)
	return data, count, err2
}

// ===
func (p *DMManage) GetOneItem(tableName string, id string) (mydb.DataItem, error) {
	return p.GetOneItemByDB(tableName, id)
}

func (p *DMManage) GetOneItemByDB(tableName string, id string) (mydb.DataItem, error) {
	filter := " and id=? "
	params := []interface{}{id}
	res, err := p.GetOneItemByDBWithFilter(tableName, filter, params)
	return res, err
}

func (p *DMManage) GetOneItemByDBWithFilter(tableName string, filter string, params []interface{}) (mydb.DataItem, error) {
	data, err := getDBOneItemFromSingleTable(p.myDbConn, tableName, filter, params, p.GetOrderBy(tableName))
	if err != nil {
		// p.mylogger.Warningf("failed to get one item, err:%v", err) //warning because it may be on purpose
		return nil, err
	}
	return data, nil
}

func (p *DMManage) GetOneItemByDBWithOrderBy(tableName string, filter string, params []interface{}, tableOrder string) (mydb.DataItem, error) {
	data, err := getDBOneItemFromSingleTable(p.myDbConn, tableName, filter, params, tableOrder)
	if err != nil {
		// p.mylogger.Warningf("failed to get one item, err:%v", err) //warning because it may be on purpose
		return nil, err
	}
	return data, nil
}

func (p *DMManage) GetOneFullItem(tableName string, id string) (mydb.DataItem, error) {
	return p.GetOneFullItemByDB(tableName, id)
}

func (p *DMManage) GetOneFullItemByDB(tableName string, id string) (mydb.DataItem, error) {
	filter := " and id=? "
	params := []interface{}{id}
	res, err := p.GetOneFullItemByDBWithFilter(tableName, filter, params, "", []interface{}{})
	return res, err
}

func (p *DMManage) GetOneFullItemByDBWithFilter(tableName string, filter string, params []interface{}, joinFilter string, joinParams []interface{}) (mydb.DataItem, error) {
	data, err := getDBJoinOneItem(p.myDbConn, tableName,
		filter, params, p.GetJoinSelect(tableName), p.GetTableJoin(tableName, ""),
		joinFilter, joinParams)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *DMManage) CreateOneItem(data mydb.DataItem) (mydb.DataItem, error) {
	if err := p.IsCreateItemValid(data); err != nil {
		return nil, fmt.Errorf("create item valid failed:%v", err)
	}
	item, err2 := createDBOneItem(p.myDbConn, data)
	if err2 != nil {
		return nil, fmt.Errorf("create item err:%v", err2)
	}
	return item, nil
}

func (p *DMManage) UpdateOneItem(data mydb.DataItem, id string) (mydb.DataItem, error) {
	if err := p.IsUpdateItemValid(data); err != nil {
		return nil, fmt.Errorf("update item valid failed:%v", err)
	}
	item, err2 := updateDBOneItem(p.myDbConn, data, id)
	if err2 != nil {
		return nil, fmt.Errorf("update item err:%v", err2)
	}
	return item, nil
}

func (p *DMManage) UpdateOneItemByForce(tableName string, id string, fields map[string]interface{}) error {
	return p.UpdateOneItemByForceWithFilter(tableName, "id=?", []interface{}{id}, fields)
}

func (p *DMManage) UpdateOneItemFlag(tableName string, id string, fields []string, flags []bool) error {
	if len(fields) != len(flags) {
		return errors.New("unmatched fileds and values")
	}
	mm := map[string]interface{}{}
	for index, field := range fields {
		mm[field] = flags[index]
	}
	return p.UpdateOneItemByForceWithFilter(tableName, " id=?", []interface{}{id}, mm)
}

func (p *DMManage) UpdateItemFieldsByForce(tableName string, whereFilter string, tableParams []interface{},
	fields []string, values []interface{}) error {
	if len(fields) != len(values) {
		return errors.New("unmatched fileds and values")
	}
	mm := map[string]interface{}{}
	for index, field := range fields {
		mm[field] = values[index]
	}
	return p.UpdateOneItemByForceWithFilter(tableName, whereFilter, tableParams, mm)
}

func (p *DMManage) UpdateOneItemByForceWithFilter(tableName string, whereFilter string, tableParams []interface{}, fields map[string]interface{}) error {
	return updateDBItemsFieldsWithFilter(p.myDbConn, tableName, whereFilter, tableParams, fields)
}

func (p *DMManage) AddOrUpdateOneItem(data mydb.DataItem,
	whereFilter string, params []interface{}) (mydb.DataItem, error) {
	var bi mydb.DataItem
	var err error
	//this is a tmp fix
	tableFilter := ""
	if strings.TrimSpace(whereFilter) != "" {
		tableFilter = " and " + whereFilter //become a tableFilter
	}
	data2, err2 := p.GetOneItemByDBWithFilter(data.TableName(), tableFilter, params)
	if err2 != nil || data2.GetID() == "" {
		if err3 := p.IsCreateItemValid(data); err3 != nil {
			return nil, fmt.Errorf("create item valid failed:%v", err3)
		}
		bi, err = createDBOneItem(p.myDbConn, data)
		if err != nil {
			return nil, fmt.Errorf("create item err:%v", err)
		}
		bi, err = p.GetOneItem(data.TableName(), bi.GetID())
	} else {
		if err3 := p.IsUpdateItemValid(data); err3 != nil {
			return nil, fmt.Errorf("update item valid failed:%v", err3)
		}
		_, err = updateDBOneItemWithFilter(p.myDbConn, data, whereFilter, params)
		if err != nil {
			return nil, err
		}
		bi, err = p.GetOneItem(data.TableName(), data2.GetID())
	}
	if err != nil {
		return nil, err
	}

	return bi, nil
}

func (p *DMManage) UpdateOneItemWithFilter(data mydb.DataItem,
	whereFilter string, params []interface{}) (mydb.DataItem, error) {
	var bi mydb.DataItem
	var err error
	if err = p.IsUpdateItemValid(data); err != nil {
		return nil, fmt.Errorf("update item valid failed:%v", err)
	}
	bi, err = updateDBOneItemWithFilter(p.myDbConn, data, whereFilter, params)
	return bi, nil
}

func (p *DMManage) SaveOneItem(data mydb.DataItem) (mydb.DataItem, error) {
	if data.GetID() == "" {
		if err := p.IsCreateItemValid(data); err != nil {
			return nil, fmt.Errorf("create item valid failed:%v", err)
		}
	} else {
		if err := p.IsUpdateItemValid(data); err != nil {
			return nil, fmt.Errorf("update item valid failed:%v", err)
		}
	}
	bi, err := SyncDBOneItem(p.myDbConn, data, data.GetID())
	if err != nil {
		return nil, fmt.Errorf("sync item err:%v", err)
	}
	return bi, nil
}

func (p *DMManage) DeleteOneItem(tableName string, id string) error {
	return p.DeleteItemWithFilter(tableName, "id=? ", []interface{}{id})
}

func (p *DMManage) DeleteItemWithFilter(tableName string, whereFilter string, params []interface{}) error {
	err := deleteDBItem(p.myDbConn, tableName, whereFilter, params)
	if err != nil {
		return err
	}
	return nil
}

// sync
func (p *DMManage) GetListByDBForSync(tableName string, param mydb.RequestParam,
	tableFilter string, tableParams []interface{}, orderBy string) ([]mydb.DataItem, error) {
	data, err := getDBSingleTableListForSync(p.myDbConn,
		tableName, param, p.GetTimeField(tableName), tableFilter, tableParams,
		orderBy, p.GetFieldMap(tableName))
	return data, err
}

func (p *DMManage) GetOneItemByDBForSync(tableName string, id string) (mydb.DataItem, error) {
	filter := " and id=? "
	params := []interface{}{id}
	data, err := getDBOneItemFromSingleTableForSync(p.myDbConn, tableName, filter, params, "")
	return data, err
}

func (p *DMManage) GetOneItemByDBWithFilterForSync(tableName string, filter string, params []interface{}, orderBy string) (mydb.DataItem, error) {
	if orderBy == "" {
		orderBy = " order by updated_at desc "
	}
	data, err := getDBOneItemFromSingleTableForSync(p.myDbConn, tableName, filter, params, orderBy)
	return data, err
}

func (p *DMManage) GetOneFullItemByDBForSync(tableName string, id string) (mydb.DataItem, error) {
	filter := " and id=? "
	params := []interface{}{id}
	data, err := getDBJoinOneItemForSync(p.myDbConn, tableName,
		filter, params,
		p.GetJoinSelect(tableName), p.GetTableJoin(tableName, ""),
		"", []interface{}{})
	return data, err
}
