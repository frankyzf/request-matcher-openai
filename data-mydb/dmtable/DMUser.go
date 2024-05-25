package dmtable

import (
	"errors"
	"request-matcher-openai/data-mydb/mydb"
)

type UserAdapter struct {
}

func (p *UserAdapter) GetName() string {
	return "user"
}

func (p *UserAdapter) CacheName() string {
	return ""
}

func (p *UserAdapter) FullItemCacheName() string {
	return ""
}

func (p *UserAdapter) GetDataItem() mydb.DataItem {
	return mydb.User{}
}

func (p *UserAdapter) GetFullDataItem() mydb.DataItem {
	return mydb.UserShort{}
}

func (p *UserAdapter) GetFieldMap() map[string]string {
	return map[string]string{}
}

func (p *UserAdapter) GetJoinSelect() string {
	return `select user.* `
}

func (p *UserAdapter) GetTableJoin() string {
	return ` where user.deleted_at is null `
}

func (p *UserAdapter) UserFieldName() string {
	return ""
}

func (p *UserAdapter) GetTimeField() string {
	return "created_at"
}

func (p *UserAdapter) GetOrderBy() string {
	return `  order by created_at desc `
}

func (p *UserAdapter) GetFullItemOrderBy() string {
	return ` order by user.created_at desc `
}

func (p *UserAdapter) IsCreateItemValid(data interface{}) error {
	if _, ok := data.(mydb.User); ok {
		return nil
	}
	return errors.New("failed to convert to User")
}

func (p *UserAdapter) IsUpdateItemValid(data interface{}) error {
	return nil
}
