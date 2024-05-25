package dmtable

import (
	"errors"
	"request-matcher-openai/data-mydb/mydb"
)

type ProjectAdapter struct {
}

func (p *ProjectAdapter) GetName() string {
	return "project"
}

func (p *ProjectAdapter) CacheName() string {
	return ""
}

func (p *ProjectAdapter) FullItemCacheName() string {
	return ""
}

func (p *ProjectAdapter) GetDataItem() mydb.DataItem {
	return mydb.Project{}
}

func (p *ProjectAdapter) GetFullDataItem() mydb.DataItem {
	return mydb.ProjectShort{}
}

func (p *ProjectAdapter) GetFieldMap() map[string]string {
	return map[string]string{}
}

func (p *ProjectAdapter) GetJoinSelect() string {
	return `select project.* `
}

func (p *ProjectAdapter) GetTableJoin() string {
	return ` where project.deleted_at is null  `
}

func (p *ProjectAdapter) UserFieldName() string {
	return ""
}

func (p *ProjectAdapter) GetTimeField() string {
	return "created_at"
}

func (p *ProjectAdapter) GetOrderBy() string {
	return " order by created_at desc "
}

func (p *ProjectAdapter) GetFullItemOrderBy() string {
	return " order by project.created_at desc "
}

func (p *ProjectAdapter) IsCreateItemValid(data interface{}) error {
	if _, ok := data.(mydb.Project); ok {
		return nil
	}
	return errors.New("failed to convert to Project")
}

func (p *ProjectAdapter) IsUpdateItemValid(data interface{}) error {
	return nil
}
