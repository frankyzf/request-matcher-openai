package db

import (
	"fmt"
	"net/url"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var MysqlHost, MysqlPort, MysqlUser, MysqlPassword, MysqlDatabase string
var MysqlUrl string

func LoadAndSetupMysql(vp *viper.Viper) *gorm.DB {
	MysqlHost = "localhost"
	if vp.IsSet("mysql.host") {
		MysqlHost = vp.GetString("mysql.host")
	}

	MysqlPort = "3306"
	if vp.IsSet("mysql.port") {
		MysqlPort = vp.GetString("mysql.port")
	}

	MysqlUser = "test"
	if vp.IsSet("mysql.user") {
		MysqlUser = vp.GetString("mysql.user")
	}

	MysqlPassword = "test123!"
	if vp.IsSet("mysql.password") {
		MysqlPassword = vp.GetString("mysql.password")
	}

	MysqlDatabase = "test"
	if vp.IsSet("mysql.database") {
		MysqlDatabase = vp.GetString("mysql.database")
	}

	timezone := "'Asia/Shanghai'"

	MysqlUrl = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?loc=Local&parseTime=true&time_zone=%s", url.QueryEscape(MysqlUser),
		MysqlPassword, MysqlHost, MysqlPort, MysqlDatabase, url.QueryEscape(timezone))
	mysqldb := InitOrm("mysql", MysqlUrl)
	// mysqldb.LogMode(true)
	return mysqldb
}
