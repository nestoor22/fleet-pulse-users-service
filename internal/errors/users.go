package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ErrUserNotFound = errors.New("user not found")
var ErrEmailAlreadyExists = errors.New("user with such email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")

func HandleUserErrors(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrEmailAlreadyExists):
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, ErrUserNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, ErrInvalidCredentials):
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
	}
}
