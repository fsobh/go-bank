package api

import (
	db "github.com/fsobh/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server define http server here 1.
type Server struct {
	store  db.Store // this used to be a pointer. we removed it after creating the querier.go interface for mocking the DB. Interfaces cannot be pointers
	router *gin.Engine
}

// NewServer Creates a new HTTP server instance and setup routing 2.
func NewServer(store db.Store) *Server {

	server := &Server{store: store}
	router := gin.Default()

	// get the validator engine for gin (as a pointer)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		//name, validator callback
		v.RegisterValidation("currency", validCurrency)
	}

	//add routes to router

	// defined in account.go
	router.POST("/accounts", server.createAccount) // POST
	router.GET("/accounts/:id", server.getAccount) // GET id is a URL params
	router.GET("/accounts", server.listAccount)    // this implements pagination for fetching all accounts in a range

	router.POST("/transfers", server.createTransfer) // this implements pagination for fetching all accounts in a range

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
