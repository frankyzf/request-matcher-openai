package export

import "request-matcher-openai/data-mydb/mydb"

type ConversionInterface interface {
	ParseOneItem(buf string) (mydb.DataItem, error)
	MarshalOneItem(item interface{}) (string, error)
	ParseOneFullItem(buf string) (mydb.DataItem, error)
	MarshalOneFullItem(item interface{}) (string, error)
	ParseItems(buf string) ([]mydb.DataItem, error)
	MarshalItems(data []mydb.DataItem) (string, error)
	ParseFullItems(buf string) ([]mydb.DataItem, error)
	MarshalFullItems(data []mydb.DataItem) (string, error)
	ConvertItems(data interface{}) ([]mydb.DataItem, error)
	ConvertFullItems(data interface{}) ([]mydb.DataItem, error)
}
