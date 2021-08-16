package db

import (
	"database/sql"
	"fmt"
	"log"
	"tgm/config"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func DbManager() *sql.DB {
	configuration := config.GetConfig()

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", configuration.DB_USERNAME, configuration.DB_PASSWORD, configuration.DB_HOST, configuration.DB_PORT, configuration.DB_NAME)

	db, err := sql.Open("mysql", connStr)

	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(configuration.DB_MAXIDLECONNS)
	db.SetMaxOpenConns(configuration.DB_MAXOPENCONNS)

	if err = db.Ping(); err != nil {
		panic(err)
	}
	return db
}
