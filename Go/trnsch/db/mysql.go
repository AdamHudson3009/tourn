package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() {
	var err error

	dsn := "root:MySQL=1@tcp(127.0.0.1:3306)/league"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB ping error:", err)
	}

	log.Println("Connected to MySQL database")
}

func ChangDB(sql string, params []interface{}) (int, error) {
	stmt, err := DB.Prepare(sql)
	if err != nil {
		return 0, err
	}

	args := make([]interface{}, len(params))
	for i, v := range params {
		args[i] = v
	}

	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastID), nil
}
