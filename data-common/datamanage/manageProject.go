package datamanage

import (
	"request-matcher-openai/data-mydb/dmtable"

	"request-matcher-openai/data-mydb/mydb"
)

func (p *DataManager) GetProjectList(param mydb.RequestParam, bMaskPassword bool) ([]mydb.Project, int, error) {
	tableFilter := ""
	tableParams := []interface{}{}
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

func (p *DataManager) GetProjectListByDB(param mydb.RequestParam, bMaskPassword bool) ([]mydb.Project, int, error) {
	tableFilter := ""
	tableParams := []interface{}{}
	data, count, err2 := p.DMManage.GetListByDB("project", param, tableFilter, tableParams)
	if err2 != nil {
		return []mydb.Project{}, count, err2
	}
	data1, err3 := dmtable.ProjectConversion{}.ConvertItemsWithType(data)
	if err3 != nil {
		return []mydb.Project{}, count, err3
	}
	return data1, count, nil
}

func (p *DataManager) GetFullProjectListWithFilter(param mydb.RequestParam,
	tableFilter string, tableParams []interface{}) ([]mydb.ProjectShort, int, error) {
	data, count, err2 := p.DMManage.GetFullItemListByDB("project", param, tableFilter, tableParams, "", []interface{}{})
	if err2 != nil {
		return []mydb.ProjectShort{}, count, err2
	}
	data1, err3 := dmtable.ProjectConversion{}.ConvertFullItemsWithType(data)
	if err3 != nil {
		return []mydb.ProjectShort{}, count, err3
	}
	return data1, count, nil
}

func (p *DataManager) GetOneProject(id string, bMaskPassword bool) (mydb.Project, error) {
	item, err := p.DMManage.GetOneItem("project", id)
	if err != nil {
		return mydb.Project{}, err
	}
	data2, err2 := dmtable.ProjectConversion{}.ConvertOneItemWithType(item)
	if err2 != nil {
		return mydb.Project{}, err2
	}
	return data2, nil
}

func (p *DataManager) GetOneProjectByDB(id string, bMaskPassword bool) (mydb.Project, error) {
	item, err := p.DMManage.GetOneItemByDB("project", id)
	if err != nil {
		return mydb.Project{}, err
	}
	data2, err2 := dmtable.ProjectConversion{}.ConvertOneItemWithType(item)
	if err2 != nil {
		return mydb.Project{}, err2
	}
	return data2, nil
}

func (p *DataManager) GetOneProjectByDBWithFilter(tableFilter string, tableParams []interface{}, bMaskPassword bool) (mydb.Project, error) {
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
	item, err := p.DMManage.GetOneFullItem("project", id)
	if err != nil {
		return mydb.ProjectShort{}, err
	}
	return dmtable.ProjectConversion{}.ConvertOneFullItemWithType(item)
}

func (p *DataManager) GetOneFullProjectByDB(id string) (mydb.ProjectShort, error) {
	item, err := p.DMManage.GetOneFullItemByDB("project", id)
	if err != nil {
		return mydb.ProjectShort{}, err
	}
	return dmtable.ProjectConversion{}.ConvertOneFullItemWithType(item)
}

func (p *DataManager) CreateOneProject(msg mydb.Project) (mydb.Project, error) {
	// use internal method
	var err error
	msg, err = p.AddOrUpdateProjectItem(msg)
	if err != nil {
		return mydb.Project{}, err
	}
	//p.sendProjectWelcomeMessage(proeject, msg.Password)
	return msg, nil
}

func (p *DataManager) UpdateOneProject(id string, msg mydb.Project) (mydb.Project, error) {
	// use internal method
	var err error
	msg, err = p.AddOrUpdateProjectItem(msg)
	if err != nil {
		return mydb.Project{}, err
	}
	return msg, nil
}

func (p *DataManager) AddOrUpdateProjectItem(proeject mydb.Project) (mydb.Project, error) {
	// use internal method
	item, err2 := p.DMManage.AddOrUpdateOneItem(proeject, "id=?", []interface{}{proeject.ID})
	if err2 != nil {
		return mydb.Project{}, err2
	}
	data, err2 := dmtable.ProjectConversion{}.ConvertOneItemWithType(item)
	if err2 != nil {
		return mydb.Project{}, err2
	}

	return p.GetOneProjectByDB(data.ID, true)
}

func (p *DataManager) DeleteOneProject(id string) error {
	err := p.DMManage.DeleteOneItem("project", id)
	return err
}
