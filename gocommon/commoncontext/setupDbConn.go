package commoncontext

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/gorm"
	"request-matcher-openai/gocommon/db"
)

func setupDatabaseConnection(vp *viper.Viper) (*gorm.DB, error) {
	Dialect := vp.GetString("dialect")
	if Dialect != "" {
		var dbconn *gorm.DB
		// if Dialect == "postgres" {
		// 	dbconn = db.LoadAndSetupPostgres(vp)
		// } else
		if Dialect == "mysql" {
			dbconn = db.LoadAndSetupMysql(vp)
		} else {
			return nil, fmt.Errorf("unknown dialect: %v", Dialect)
		}
		if dbconn == nil {
			GetInstance().MyLogger.Infof("failed to setup database")
			return nil, errors.New("null database")
		}
		return dbconn, nil
	} else {
		GetInstance().MyLogger.Infof("skip to setup database")
		return nil, nil
	}
}
