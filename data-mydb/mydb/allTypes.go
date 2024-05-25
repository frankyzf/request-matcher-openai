package mydb

type DataTypeParameter struct {
	DataItem DataItem
}

// this is the table index for all sql tables, and it should be maintained when new table is insert
var allDataTypes = map[string]DataTypeParameter{ //should never change
	"user": DataTypeParameter{
		DataItem: User{},
	},
	"account_user": DataTypeParameter{
		DataItem: Account{},
	},
}

func IsValidDataType(tableName string) bool {
	_, ok := allDataTypes[tableName]
	return ok
}

func GetAllDataTypes() map[string]DataItem {
	mm := map[string]DataItem{}
	for tableName, dtParameter := range allDataTypes {
		mm[tableName] = dtParameter.DataItem
	}
	return mm
}

func GetAllDataTypeParameters() map[string]DataTypeParameter {
	return allDataTypes
}
