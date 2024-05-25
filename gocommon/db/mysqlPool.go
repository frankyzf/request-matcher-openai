package db

import (
	_ "github.com/go-sql-driver/mysql"
)

// var Mysql *sql.DB

// func SetupMyql(url string) {
// 	Mysql = InitMysql(url)
// }

// func InitMysql(url string) *sql.DB {
// 	fmt.Printf("mysql pool init with url: %s\n", url)
// 	db, _ := sql.Open("mysql", url)
// 	db.SetMaxOpenConns(100)
// 	db.SetMaxIdleConns(100)
// 	db.SetConnMaxLifetime(time.Minute)
// 	db.Ping()
// 	return db
// }

// func WriteMysql(p *sql.DB, sqlStr string) int64 {
// 	rs, err := p.Exec(sqlStr)
// 	if err != nil {
// 		fmt.Printf("write sqlStr:%s encounter a error: %v\n", sqlStr, err)
// 		return 0
// 	}
// 	rowCount, err := rs.RowsAffected()
// 	if err != nil {
// 		fmt.Printf("get error in rowsAffected: %v\n", err)
// 		return 0
// 	}
// 	return rowCount
// }
