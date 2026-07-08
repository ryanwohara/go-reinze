package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Db() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_0900_ai_ci", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASS"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_DATABASE"))

	var err error
	var database *sql.DB

	// sql.Open only validates the DSN; Ping is what actually dials,
	// so retry until the database is reachable.
	for {
		database, err = sql.Open("mysql", dsn)

		if err == nil {
			err = database.Ping()
		}

		if err == nil {
			break
		}

		if database != nil {
			database.Close()
		}

		fmt.Println("Error connecting to database:", err)
		time.Sleep(time.Second)
	}

	database.SetConnMaxLifetime(time.Second * 5)
	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(10)

	fmt.Printf("Database connection established at %s:%s\n", os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"))

	return database
}
