package dmmodule

import (
	"request-matcher-openai/data-mydb/mydb"
)

type DMAdapterInterface interface {
	GetName() string
	CacheName() string
	FullItemCacheName() string
	GetDataItem() mydb.DataItem
	GetFullDataItem() mydb.DataItem
	GetFieldMap() map[string]string
	GetJoinSelect() string
	GetTableJoin() string
	UserFieldName() string
	GetTimeField() string
	GetOrderBy() string
	GetFullItemOrderBy() string
	IsCreateItemValid(data interface{}) error
	IsUpdateItemValid(data interface{}) error
}
