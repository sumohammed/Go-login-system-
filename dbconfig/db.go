package dbconfig

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
	"log"
)

var DB *sql.DB

func InitDB(data string) {
	var err error
	DB, err = sql.Open("mysql", data)
	if err != nil {
		log.Panic(err)
	}

	if err = DB.Ping(); err != nil {
		log.Panic(err)
	}
}