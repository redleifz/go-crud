package configs

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
)

var DB *sql.DB

// func EnvMySql() string {
// 	// You should return the complete MySQL DSN as a string.
// 	return "username:password@tcp(localhost:3306)/employee"
// }

func ConnectDB() *sql.DB {
	var err error

	DB, err = sql.Open("mysql", EnvMySql())
	if err != nil {
		log.Fatal(err)
	}

	// Test the connection to the database
	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Connected to DB!")
	return DB
}
