package main

//7
import (
	"database/sql"
	"github.com/fsobh/simplebank/api"
	db "github.com/fsobh/simplebank/db/sqlc"
	"github.com/fsobh/simplebank/util"
	_ "github.com/lib/pq" // VERY IMPORTANT FOR DB TO WORK ON SERVER
	"log"
)

// Main entry point into server
func main() {

	//load env file
	config, err := util.LoadConfig(".") // pass in the path relative to this file
	if err != nil {
		log.Fatal("Cannot load env variables", err)
	}

	// Create an instance to our DB connection
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	//error handling
	if err != nil {
		log.Fatal("Can not connect to db:", err)
	}

	//Declare and initialize a new database store instance by passing in the connection instance to the database
	store := db.NewStore(conn)
	//Declare and initialize an instance of our server by passing in the database store instance
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start server")
	}
}
