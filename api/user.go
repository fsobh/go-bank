package api

import (
	"errors"
	db "github.com/fsobh/simplebank/db/sqlc"
	"github.com/fsobh/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"time"
)

// These are provided by the default validator package :
// alphanum means no special chars
// min=6 means at least 6 chars
// email is for verifying value is email addy
type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type CreateUserResponse struct {
	Username         string    `json:"username"`
	FullName         string    `json:"full_name"`
	Email            string    `json:"email"`
	PasswordChangeAt time.Time `json:"password_change_at"`
	CreatedAt        time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {

	var req createUserRequest

	// We validate the input json with what we said we'd expect in the struct (using the context)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//If the validation fails, return a 400 error code
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)

	// make sure there's no errors
	if err != nil {
		// try to convert it to pq error (to handle violation of constraints)
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			// Only check for unique violations since user has no fk.
			//This will check if a user with the same username/email is already signed up
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rsp := CreateUserResponse{
		Username:         user.Username,
		FullName:         user.FullName,
		Email:            user.Email,
		PasswordChangeAt: user.PasswordChangeAt,
		CreatedAt:        user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, rsp)

}
