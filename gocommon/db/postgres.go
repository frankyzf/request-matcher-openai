package db

// var PostgresHost, PostgresPort, PostgresUser, PostgresPassword, PostgresDatabase string
// var PostgresUrl string

// func LoadAndSetupPostgres(vp *viper.Viper) *gorm.DB {
// 	PostgresHost = "localhost"
// 	if vp.IsSet("postgres.host") {
// 		PostgresHost = vp.GetString("postgres.host")
// 	}

// 	PostgresPort = "5432"
// 	if vp.IsSet("postgres.port") {
// 		PostgresPort = vp.GetString("postgres.port")
// 	}

// 	PostgresUser = "postgres"
// 	if vp.IsSet("postgres.user") {
// 		PostgresUser = vp.GetString("postgres.user")
// 	}

// 	PostgresPassword = "postgres"
// 	if vp.IsSet("postgres.password") {
// 		PostgresPassword = vp.GetString("postgres.password")
// 	}

// 	PostgresDatabase = "test"
// 	if vp.IsSet("postgres.database") {
// 		PostgresDatabase = vp.GetString("postgres.database")
// 	}

// 	PostgresUrl = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable ",
// 		PostgresHost, PostgresPort, PostgresUser, PostgresDatabase, PostgresPassword)

// 	postgresdb := InitOrm("postgres", PostgresUrl)
// 	postgresdb.LogMode(true)
// 	return postgresdb
// }
