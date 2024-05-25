package datamanage

import (
	"request-matcher-openai/data-mydb/dmtable"
	"request-matcher-openai/data-mydb/mydb"
)

func (p *DataManager) GetQualificationMatchList(param mydb.RequestParam, tableFilter string, tableParams []interface{}) ([]mydb.QualificationMatch, int, error) {
	data, count, err2 := p.DMManage.GetList("qualification_match", param, tableFilter, tableParams)
	if err2 != nil {
		return []mydb.QualificationMatch{}, count, err2
	}
	data1, err3 := dmtable.QualificationMatchConversion{}.ConvertItemsWithType(data)
	if err3 != nil {
		return []mydb.QualificationMatch{}, count, err3
	}
	return data1, count, nil
}

func (p *DataManager) GetFullQualificationMatchList(param mydb.RequestParam, tableFilter string, tableParams []interface{}) ([]mydb.QualificationMatchShort, int, error) {
	data, count, err2 := p.DMManage.GetFullItemList("qualification_match", param, tableFilter, tableParams, "", []interface{}{})
	if err2 != nil {
		return []mydb.QualificationMatchShort{}, count, err2
	}
	data1, err3 := dmtable.QualificationMatchConversion{}.ConvertFullItemsWithType(data)
	if err3 != nil {
		return []mydb.QualificationMatchShort{}, count, err3
	}
	return data1, count, nil
}

func (p *DataManager) SplitGetQualificationMatchList(param mydb.RequestParam, tableFilter string, tableParams []interface{}) ([]mydb.QualificationMatch, error) {
	data, err2 := p.DMManage.SplitGetListByDB("qualification_match", param, tableFilter, tableParams)
	if err2 != nil {
		return []mydb.QualificationMatch{}, err2
	}
	data1, err3 := dmtable.QualificationMatchConversion{}.ConvertItemsWithType(data)
	if err3 != nil {
		return []mydb.QualificationMatch{}, err3
	}
	return data1, nil
}

func (p *DataManager) GetOneQualificationMatch(id string) (mydb.QualificationMatch, error) {
	item, err := p.DMManage.GetOneItem("qualification_match", id)
	if err != nil {
		return mydb.QualificationMatch{}, err
	}
	return dmtable.QualificationMatchConversion{}.ConvertOneItemWithType(item)
}

func (p *DataManager) GetOneQualificationMatchByDBWithFilter(tableFilter string, tableParams []interface{}) (mydb.QualificationMatch, error) {
	item, err := p.DMManage.GetOneItemByDBWithFilter("qualification_match", tableFilter, tableParams)
	if err != nil {
		return mydb.QualificationMatch{}, err
	}
	data2, err2 := dmtable.QualificationMatchConversion{}.ConvertOneItemWithType(item)
	if err2 != nil {
		return mydb.QualificationMatch{}, err2
	}
	return data2, nil
}

func (p *DataManager) GetOneFullQualificationMatch(id string) (mydb.QualificationMatchShort, error) {
	item, err2 := p.DMManage.GetOneFullItem("qualification_match", id)
	if err2 != nil {
		return mydb.QualificationMatchShort{}, err2
	}
	data, err := dmtable.QualificationMatchConversion{}.ConvertOneFullItemWithType(item)

	return data, err
}

func (p *DataManager) UpdateOneQualificationMatch(id string, data mydb.QualificationMatch) (mydb.QualificationMatch, error) {
	data.ID = id
	item, err2 := p.DMManage.UpdateOneItem(data, id)
	if err2 != nil {
		return mydb.QualificationMatch{}, err2
	}
	return dmtable.QualificationMatchConversion{}.ConvertOneItemWithType(item)
}

func (p *DataManager) DeleteOneQualificationMatch(id string) error {
	err := p.DMManage.DeleteOneItem("qualification_match", id)
	return err
}
