package main

import (
	"database/sql"
	"os"

	"github.com/dascr/dascr-board/api"
	"github.com/dascr/dascr-board/config"
	"github.com/dascr/dascr-board/database"
	"github.com/dascr/dascr-board/logger"
)

var (
	db  *sql.DB
	err error
	// Debug will hold global debug flag
	Debug bool = false
)

func main() {
	// Generate uploads directory
	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
		logger.Panicf("Unable to create uploads directory: %+v", err)
	}
	// Setup DB
	dbconfig := &config.DBConfig{
		Driver:   "sqlite3",
		Filename: "./dascr.db",
	}
	if db, err = database.SetupDB(dbconfig); err != nil {
		logger.Panicf("Unable to create database: %+v", err)
	}

	// API Config
	APIConfig := &config.APIConfig{
		IP:   config.MustGet("API_IP"),
		Port: config.MustGet("API_PORT"),
	}

	// Setup API
	a := api.SetupAPI(db, APIConfig)
	a.Start()
}
