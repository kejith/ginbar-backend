//main.go

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"ginbar/fiberapi"
	"ginbar/mysql/db"

	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

//var secret = "IX~|xTE@4*v@e95sLll4g`#6G288be"

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
	// Load Configs from .env File
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// Database Config
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

	// Create DSN String for Opening a Connection to the Database
	dsn := fmt.Sprintf(
		"%s:%s@/%s?parseTime=true",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DatabaseName,
	)

	// Open and Store Connection to the Database
	dbConnection, err = sql.Open(dbConfig.Driver, dsn)
	if err != nil {
		log.Fatal("Can't connect to mysql", err)
	}

	dbStatsJSON, err := json.Marshal(dbConnection.Stats())
	if err != nil {
		log.Fatalf("could not convert sql.dbStats to json: %v", err)
	}

	// Settings
	fmt.Printf("%s", string(dbStatsJSON))
}

func main() {
	/*
	* Stors is used for all kind of Database related functions
	 */
	store := db.NewStore(dbConnection)

	/*
	* Server handles Incoming HTTP Requests, Store the State of the Server
	* and Connections
	 */

	server, err := fiberapi.NewFiber(store)
	if err != nil {
		log.Fatal("Can't create server", err)
	}
	_ = server

	// Start Server
	// server, err := api.NewServer(store)
	// if err != nil {
	// 	log.Fatal("Can't create server", err)
	// }
	// _ = server
	// err = server.Start(":8080")
	// if err != nil {
	// 	log.Fatal("Can't start server", err)
	// }
}
