package api

import (
	db "github.com/fsobh/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server define http server here 1.
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer Creates a new HTTP server instance and setup routing 2.
func NewServer(store *db.Store) *Server {

	server := &Server{store: store}
	router := gin.Default()

	//add routes to router

	router.POST("/accounts", server.createAccount) // defined in account.go

	server.router = router
	return server

}

// Start Starts the server to run and listen on the given address 6
// We Create this because the router field in our server struct is private (cant be accessed outside of this package)
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// This is a global function that converts go errors into a key value object 5
// This will allow Gin to serialize it to JSON so that we can return it back  to the client (using context)
// gin.H is a shortcut for map interface
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}