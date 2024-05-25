package dmmodule

import (
	"errors"
	"request-matcher-openai/data-mydb/dmtable"

	"request-matcher-openai/data-mydb/mydb"
)

func getDMUserTypeAdapter(name string) (DMAdapterUserTypeInterface, error) {
	if mydb.IsUserTypeSensitive(name) == false {
		return nil, errors.New("not user type sensitive")
	}

	mylogger.Warningf("skip to construct dmitem:%v, because unknown user type:%v", name, name)
	return nil, errors.New("unknown dm user type for table:" + name)
}

func getDMStatusTypeAdapter(name string) (DMAdapterStatusTypeInterface, error) {
	mylogger.Warningf("skip to construct dmitem:%v, because unknown status type:%v", name, name)
	return nil, errors.New("unknown dm status type for table:" + name)
}

func getDMGeneralAdapter(name string) (DMAdapterInterface, error) {
	if name == "user" {
		return &dmtable.UserAdapter{}, nil
	} else if name == "account_user" {
		return &dmtable.AccountAdapter{}, nil
	} else if name == "project" {
		return &dmtable.ProjectAdapter{}, nil
	}
	mylogger.Warningf("skip to construct dmitem:%v, because unknown typename:%v", name, name)
	return nil, errors.New("unknown dm adapter and conversion:" + name)
}
