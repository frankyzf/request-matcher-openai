package datamanage

import (
	"request-matcher-openai/data-mydb/dmtable"
	"request-matcher-openai/data-mydb/mydb"
)

func (p *DataManager) GetProjectList(param mydb.RequestParam, tableFilter string, tableParams []interface{}) ([]mydb.Project, int, error) {
	data, count, err2 := p.DMManage.GetList("project", param, tableFilter, tableParams)
	if err2 != nil {
		return []mydb.Project{}, count, err2
	}
	data1, err3 := dmtable.ProjectConversion{}.ConvertItemsWithType(data)
	if err3 != nil {
		return []mydb.Project{}, count, err3
	}
	return data1, count, nil
}

func (p *DataManager) GetFullProjectList(param mydb.RequestParam, tableFilter string, tableParams []interface{}) ([]mydb.ProjectShort, int, error) {
	data, count, err2 := p.DMManage.GetFullItemList("project", param, tableFilter, tableParams, "", []interface{}{})
	if err2 != nil {
		return []mydb.ProjectShort{}, count, err2
	}
	data1, err3 := dmtable.ProjectConversion{}.ConvertFullItemsWithType(data)
	if err3 != nil {
		return []mydb.ProjectShort{}, count, err3
	}
	return data1, count, nil
}

func (p *DataManager) SplitGetProjectList(param mydb.RequestParam, tableFilter string, tableParams []interface{}) ([]mydb.Project, error) {
	data, err2 := p.DMManage.SplitGetListByDB("project", param, tableFilter, tableParams)
	if err2 != nil {
		return []mydb.Project{}, err2
	}
	data1, err3 := dmtable.ProjectConversion{}.ConvertItemsWithType(data)
	if err3 != nil {
		return []mydb.Project{}, err3
	}
	return data1, nil
}

func (p *DataManager) GetOneProject(id string) (mydb.Project, error) {
	item, err := p.DMManage.GetOneItem("project", id)
	if err != nil {
		return mydb.Project{}, err
	}
	return dmtable.ProjectConversion{}.ConvertOneItemWithType(item)
}

func (p *DataManager) GetOneProjectByDBWithFilter(tableFilter string, tableParams []interface{}) (mydb.Project, error) {
	item, err := p.DMManage.GetOneItemByDBWithFilter("project", tableFilter, tableParams)
	if err != nil {
		return mydb.Project{}, err
	}
	data2, err2 := dmtable.ProjectConversion{}.ConvertOneItemWithType(item)
	if err2 != nil {
		return mydb.Project{}, err2
	}
	return data2, nil
}

func (p *DataManager) GetOneFullProject(id string) (mydb.ProjectShort, error) {
	item, err2 := p.DMManage.GetOneFullItem("project", id)
	if err2 != nil {
		return mydb.ProjectShort{}, err2
	}
	data, err := dmtable.ProjectConversion{}.ConvertOneFullItemWithType(item)

	return data, err
}

func (p *DataManager) CreateOneProject(req mydb.ProjectMessage) (mydb.Project, error) {
	data := mydb.ConvertMessageToProject(req)
	item, err2 := p.DMManage.CreateOneItem(data)
	if err2 != nil {
		return mydb.Project{}, err2
	}
	return dmtable.ProjectConversion{}.ConvertOneItemWithType(item)
}

func (p *DataManager) UpdateOneProject(id string, data mydb.Project) (mydb.Project, error) {
	data.ID = id
	item, err2 := p.DMManage.UpdateOneItem(data, id)
	if err2 != nil {
		return mydb.Project{}, err2
	}
	return dmtable.ProjectConversion{}.ConvertOneItemWithType(item)
}

func (p *DataManager) DeleteOneProject(id string) error {
	err := p.DMManage.DeleteOneItem("project", id)
	return err
}
