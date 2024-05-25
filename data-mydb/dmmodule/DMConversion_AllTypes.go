package dmmodule

import (
	"request-matcher-openai/data-mydb/dmtable"
	"request-matcher-openai/data-mydb/export"
)

func getConversion(name string) export.ConversionInterface {
	if name == "user" {
		return &dmtable.UserConversion{}
	} else if name == "account_user" {
		return &dmtable.AccountConversion{}
	}
	return nil
}
