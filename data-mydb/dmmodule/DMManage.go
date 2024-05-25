package dmmodule

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"request-matcher-openai/data-mydb/export"
	"request-matcher-openai/data-mydb/mydb"
	"request-matcher-openai/gocommon/commoncontext"
)

type DMManage struct {
	GeneralAdapters    map[string]DMAdapterInterface
	UserTypeAdapters   map[string]DMAdapterUserTypeInterface
	StatusTypeAdapters map[string]DMAdapterStatusTypeInterface
	Conversions        map[string]export.ConversionInterface
	myDbConn           *gorm.DB
	myRClient          *redis.Client
	mylogger           *logrus.Entry
	loggerLevel        string
}

func GetDMManage(db *gorm.DB,
	rclient *redis.Client) *DMManage {
	p := &DMManage{
		GeneralAdapters:    map[string]DMAdapterInterface{},
		UserTypeAdapters:   map[string]DMAdapterUserTypeInterface{},
		StatusTypeAdapters: map[string]DMAdapterStatusTypeInterface{},
		Conversions:        make(map[string]export.ConversionInterface),
		myDbConn:           db,
		myRClient:          rclient,
		mylogger:           commoncontext.SetupLogging("dm", "manage"),
	}
	return p
}

func (p *DMManage) Init() {
	for tableName := range mydb.GetAllDataTypes() {
		if mydb.IsUserTypeSensitive(tableName) {
			adapter, err := getDMUserTypeAdapter(tableName)
			if err != nil {
				p.mylogger.Errorf("failed to get user type adpter:%v, err:%v", tableName, err)
			} else {
				p.UserTypeAdapters[tableName] = adapter
			}
		} else if mydb.IsStatusTypeSensitive(tableName) {
			adapter, err := getDMStatusTypeAdapter(tableName)
			if err != nil {
				p.mylogger.Errorf("failed to get status type adpter:%v, err:%v", tableName, err)
			} else {
				p.StatusTypeAdapters[tableName] = adapter
			}
		} else {
			adapter, err := getDMGeneralAdapter(tableName)
			if err != nil {
				p.mylogger.Errorf("failed to get user type adpter:%v, err:%v", tableName, err)
			} else {
				p.GeneralAdapters[tableName] = adapter
			}
		}
		conversion := getConversion(tableName)
		if conversion == nil {
			p.mylogger.Errorf("failed to get conversion for table:%v", tableName)
		} else {
			p.Conversions[tableName] = conversion
		}
	}
	p.mylogger.Infof("DMManage Init Finished")
}

func (p *DMManage) GetConversion(tableName string) export.ConversionInterface {
	if item, ok := p.Conversions[tableName]; ok {
		return item
	} else {
		p.mylogger.Errorf("failed GetConversion:%v", tableName)
		return nil
	}
}

func (p *DMManage) GetDBConn() *gorm.DB {
	return p.myDbConn
}

func (p *DMManage) GetTimeField(tableName string) string {
	if mydb.IsUserTypeSensitive(tableName) {
		if item, ok := p.UserTypeAdapters[tableName]; ok {
			return item.GetTimeField()
		}
	} else if mydb.IsStatusTypeSensitive(tableName) {
		if item, ok := p.StatusTypeAdapters[tableName]; ok {
			return item.GetTimeField()
		}
	} else {
		if item, ok := p.GeneralAdapters[tableName]; ok {
			return item.GetTimeField()
		}
	}
	p.mylogger.Errorf("failed GetTimeField:%v", tableName)
	return ""
}

func (p *DMManage) GetOrderBy(tableName string) string {
	if mydb.IsUserTypeSensitive(tableName) {
		if item, ok := p.UserTypeAdapters[tableName]; ok {
			return item.GetOrderBy()
		}
	} else if mydb.IsStatusTypeSensitive(tableName) {
		if item, ok := p.StatusTypeAdapters[tableName]; ok {
			return item.GetOrderBy()
		}
	} else {
		if item, ok := p.GeneralAdapters[tableName]; ok {
			return item.GetOrderBy()
		}
	}
	p.mylogger.Errorf("failed GetOrderBy:%v", tableName)
	return ""
}

func (p *DMManage) GetFieldMap(tableName string) map[string]string {
	if mydb.IsUserTypeSensitive(tableName) {
		if item, ok := p.UserTypeAdapters[tableName]; ok {
			return item.GetFieldMap()
		}
	} else if mydb.IsStatusTypeSensitive(tableName) {
		if item, ok := p.StatusTypeAdapters[tableName]; ok {
			return item.GetFieldMap()
		}
	} else {
		if item, ok := p.GeneralAdapters[tableName]; ok {
			return item.GetFieldMap()
		}
	}
	p.mylogger.Errorf("failed GetFieldMap:%v", tableName)
	return map[string]string{}
}

func (p *DMManage) GetJoinSelect(tableName string) string {
	if mydb.IsUserTypeSensitive(tableName) {
		if item, ok := p.UserTypeAdapters[tableName]; ok {
			return item.GetJoinSelect()
		}
	} else if mydb.IsStatusTypeSensitive(tableName) {
		if item, ok := p.StatusTypeAdapters[tableName]; ok {
			return item.GetJoinSelect()
		}
	} else {
		if item, ok := p.GeneralAdapters[tableName]; ok {
			return item.GetJoinSelect()
		}
	}
	p.mylogger.Errorf("failed GetJoinSelect:%v", tableName)
	return ""
}

func (p *DMManage) GetTableJoin(tableName string, userType string) string {
	if mydb.IsUserTypeSensitive(tableName) {
		if item, ok := p.UserTypeAdapters[tableName]; ok {
			return item.GetTableJoin(userType)
		}
	} else if mydb.IsStatusTypeSensitive(tableName) {
		if item, ok := p.StatusTypeAdapters[tableName]; ok {
			return item.GetTableJoin("")
		}
	} else {
		if item, ok := p.GeneralAdapters[tableName]; ok {
			return item.GetTableJoin()
		}
	}
	p.mylogger.Errorf("failed GetTableJoin:%v", tableName)
	return ""
}

func (p *DMManage) GetFullItemOrderBy(tableName string) string {
	if mydb.IsUserTypeSensitive(tableName) {
		if item, ok := p.UserTypeAdapters[tableName]; ok {
			return item.GetFullItemOrderBy()
		}
	} else if mydb.IsStatusTypeSensitive(tableName) {
		if item, ok := p.StatusTypeAdapters[tableName]; ok {
			return item.GetFullItemOrderBy()
		}
	} else {
		if item, ok := p.GeneralAdapters[tableName]; ok {
			return item.GetFullItemOrderBy()
		}
	}
	p.mylogger.Errorf("failed GetFullItemOrderBy:%v", tableName)
	return ""
}

func (p *DMManage) IsCreateItemValid(data mydb.DataItem) error {
	tableName := data.TableName()
	if mydb.IsUserTypeSensitive(tableName) {
		if item, ok := p.UserTypeAdapters[tableName]; ok {
			return item.IsCreateItemValid(data)
		}
	} else if mydb.IsStatusTypeSensitive(tableName) {
		if item, ok := p.StatusTypeAdapters[tableName]; ok {
			return item.IsCreateItemValid(data)
		}
	} else {
		if item, ok := p.GeneralAdapters[tableName]; ok {
			return item.IsCreateItemValid(data)
		}
	}
	return fmt.Errorf("unknown type:%v", tableName)
}

func (p *DMManage) IsUpdateItemValid(data mydb.DataItem) error {
	tableName := data.TableName()
	if mydb.IsUserTypeSensitive(tableName) {
		if item, ok := p.UserTypeAdapters[tableName]; ok {
			return item.IsUpdateItemValid(data)
		}
	} else if mydb.IsStatusTypeSensitive(tableName) {
		if item, ok := p.StatusTypeAdapters[tableName]; ok {
			return item.IsUpdateItemValid(data)
		}
	} else {
		if item, ok := p.GeneralAdapters[tableName]; ok {
			return item.IsUpdateItemValid(data)
		}
	}
	return fmt.Errorf("unknown type:%v", tableName)
}

func (p *DMManage) GetOneItemForSync(tableName string, id string) (mydb.DataItem, error) {
	dataItem, err2 := p.GetOneItemByDBForSync(tableName, id)
	return dataItem, err2
}
