//main.go

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"ginbar/api"
	"ginbar/mysql/db"
)

var secret = "IX~|xTE@4*v@e95sLll4g`#6G288be"

// DatabaseConfig is a struct to hold data to open a Database connection
type DatabaseConfig struct {
	Driver       string
	User         string
	Password     string
	Port         int
	Host         string
	DatabaseName string
}

var dbConfig DatabaseConfig
var dbConnection *sql.DB

func init() {
	fmt.Println("init")
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		panic(err)
	}

	dbConfig = DatabaseConfig{
		Driver:       os.Getenv("DB_DRIVER"),
		User:         os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASSWORD"),
		Port:         port,
		Host:         os.Getenv("DB_HOST"),
		DatabaseName: os.Getenv("DB_NAME"),
	}

	dsn := fmt.Sprintf(
		"%s:%s@/%s?parseTime=true",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DatabaseName,
	)

	dbConnection, err = sql.Open(dbConfig.Driver, dsn)
	if err != nil {
		log.Fatal("Can't connect to mysql", err)
	}
}

func main() {
	store := db.NewStore(dbConnection)
	server := api.NewServer(store)

	err := server.Start(":8080")
	if err != nil {
		log.Fatal("Can't start server", err)
	}
}
