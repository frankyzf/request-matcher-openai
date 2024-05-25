package dmtable

import (
	"errors"
	"request-matcher-openai/data-mydb/mydb"
)

type AccountAdapter struct {
}

func (p *AccountAdapter) GetName() string {
	return "account_user"
}

func (p *AccountAdapter) CacheName() string {
	return ""
}

func (p *AccountAdapter) FullItemCacheName() string {
	return ""
}

func (p *AccountAdapter) GetDataItem() mydb.DataItem {
	return mydb.Account{}
}

func (p *AccountAdapter) GetFullDataItem() mydb.DataItem {
	return mydb.AccountShort{}
}

func (p *AccountAdapter) GetFieldMap() map[string]string {
	return map[string]string{
		"email": "email",
		"name":  "name",
		"phone": "phone",
	}
}

func (p *AccountAdapter) GetJoinSelect() string {
	return ` select account_user.* `
}

func (p *AccountAdapter) GetTableJoin() string {
	return ` where user.deleted_at is null `
}

func (p *AccountAdapter) UserFieldName() string {
	return ""
}

func (p *AccountAdapter) GetTimeField() string {
	return "updated_at"
}

func (p *AccountAdapter) GetOrderBy() string {
	return " order by created_at desc "
}

func (p *AccountAdapter) GetFullItemOrderBy() string {
	return " order by account_user.created_at desc "
}

func (p *AccountAdapter) IsCreateItemValid(data interface{}) error {
	if _, ok := data.(mydb.Account); ok {
		return nil
	}
	return errors.New("failed to convert to Account")
}

func (p *AccountAdapter) IsUpdateItemValid(data interface{}) error {
	return nil
}
