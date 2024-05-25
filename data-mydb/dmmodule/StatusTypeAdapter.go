package dmmodule

import (
	"request-matcher-openai/data-mydb/mydb"
)

type DMAdapterStatusTypeInterface interface {
	GetName() string
	GetDataItem() mydb.DataItem
	GetFullDataItem() mydb.DataItem
	GetFieldMap() map[string]string
	GetJoinSelect() string
	GetTableJoin(status string) string
	UserFieldName() string
	GetTimeField() string
	GetOrderBy() string
	GetFullItemOrderBy() string
	IsCreateItemValid(data interface{}) error
	IsUpdateItemValid(data interface{}) error
}
