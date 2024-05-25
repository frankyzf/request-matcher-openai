package dmtable

import (
	"errors"
	"request-matcher-openai/data-mydb/mydb"
)

type QualificationMatchAdapter struct {
}

func (p *QualificationMatchAdapter) GetName() string {
	return "qualification_match"
}

func (p *QualificationMatchAdapter) CacheName() string {
	return ""
}

func (p *QualificationMatchAdapter) FullItemCacheName() string {
	return ""
}

func (p *QualificationMatchAdapter) GetDataItem() mydb.DataItem {
	return mydb.QualificationMatch{}
}

func (p *QualificationMatchAdapter) GetFullDataItem() mydb.DataItem {
	return mydb.QualificationMatchShort{}
}

func (p *QualificationMatchAdapter) GetFieldMap() map[string]string {
	return map[string]string{}
}

func (p *QualificationMatchAdapter) GetJoinSelect() string {
	return `select qualification_match.* `
}

func (p *QualificationMatchAdapter) GetTableJoin() string {
	return ` where qualification_match.deleted_at is null  `
}

func (p *QualificationMatchAdapter) UserFieldName() string {
	return ""
}

func (p *QualificationMatchAdapter) GetTimeField() string {
	return "created_at"
}

func (p *QualificationMatchAdapter) GetOrderBy() string {
	return " order by created_at desc "
}

func (p *QualificationMatchAdapter) GetFullItemOrderBy() string {
	return " order by qualification_match.created_at desc "
}

func (p *QualificationMatchAdapter) IsCreateItemValid(data interface{}) error {
	if _, ok := data.(mydb.QualificationMatch); ok {
		return nil
	}
	return errors.New("failed to convert to QualificationMatch")
}

func (p *QualificationMatchAdapter) IsUpdateItemValid(data interface{}) error {
	return nil
}
