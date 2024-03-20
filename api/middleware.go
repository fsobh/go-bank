package api

import (
	"errors"
	"fmt"
	"github.com/fsobh/simplebank/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// Higher order function for the middleware
func authMiddleWare(token token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the Value of the header
		authorizationHeader := c.GetHeader(authorizationHeaderKey)

		//Check if a value was provided
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Split the Header string using space delimiter (separate 'Bearer' from <token>)
		fields := strings.Fields(authorizationHeader)

		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Check the Authorization type of the token (first portion of token)
		authorizationType := strings.ToLower(fields[0])

		// Check if the auth type matches type bearer
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]

		// Tell our token maker to verify the token
		payload, err := token.VerifyToken(accessToken)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Set the payload in the context so the next handler can access it
		c.Set(authorizationPayloadKey, payload)
		c.Next()

	}

}
