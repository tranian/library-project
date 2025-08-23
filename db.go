package main

import (
	"database/sql"
	"fmt"
	"time"
	"log"
	_"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
	var err error
	dsn := "library_dbuser:library_securepassword@tcp(127.0.0.1:3306)/library_db?parseTime=true"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("error opening database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	fmt.Println("Connected to MariaDB")
}
