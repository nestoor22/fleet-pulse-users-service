package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ErrInvalidToken = errors.New("invalid token")
var ErrExpiredToken = errors.New("user with such email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")

func HandleAuthErrors(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidToken):
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, ErrExpiredToken):
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, ErrInvalidCredentials):
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
	}
}
