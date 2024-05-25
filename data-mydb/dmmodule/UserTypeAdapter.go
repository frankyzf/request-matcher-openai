package dmmodule

import (
	"request-matcher-openai/data-mydb/mydb"
)

type DMAdapterUserTypeInterface interface {
	GetName() string
	GetDataItem() mydb.DataItem
	GetFullDataItem() mydb.DataItem
	GetFieldMap() map[string]string
	GetJoinSelect() string
	GetTableJoin(userType string) string
	UserFieldName() string
	GetTimeField() string
	GetOrderBy() string
	GetFullItemOrderBy() string
	IsCreateItemValid(data interface{}) error
	IsUpdateItemValid(data interface{}) error
}
