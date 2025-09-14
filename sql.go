package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

func Db() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_0900_ai_ci", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASS"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_DATABASE"))

	var err error
	var database *sql.DB

	for ok := true; ok; err = nil {
		database, err = sql.Open("mysql", dsn)

		if err != nil {
			fmt.Println("Error connecting to database:", err)
			time.Sleep(time.Second)
			continue
		}

		break
	}

	database.SetConnMaxLifetime(time.Second * 5)
	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(10)

	println("Database connection established at %s:%s", os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"))

	return database
}
