package main

//7
import (
	"database/sql"
	"github.com/fsobh/simplebank/api"
	db "github.com/fsobh/simplebank/db/sqlc"
	_ "github.com/lib/pq" // VERY IMPORTANT FOR DB TO WORK ON SERVER
	"log"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

// Main entry point into server
func main() {

	// Create an instance to our DB connection
	conn, err := sql.Open(dbDriver, dbSource)

	//error handling
	if err != nil {
		log.Fatal("Can not connect to db:", err)
	}

	//Declare and initialize a new database store instance by passing in the connection instance to the database
	store := db.NewStore(conn)
	//Declare and initialize an instance of our server by passing in the database store instance
	server := api.NewServer(store)

	err = server.Start(serverAddress)

	if err != nil {
		log.Fatal("Cannot start server")
	}
}
