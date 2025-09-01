package api

import (
	"fleet-pulse-users-service/internal/errors"
	"fleet-pulse-users-service/internal/schemas"
	"fleet-pulse-users-service/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoginUserHandler godoc
// @Summary Login User
// @Description Login User
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body schemas.LoginUserRequest true "Login Data"
// @Success 200 {object} schemas.LoginResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 409 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /v1/login [post]
func LoginUserHandler(authService *services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req schemas.LoginUserRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		accessToken, refreshToken, err := authService.LoginUser(req)
		fmt.Print(err)
		if err != nil {
			errors.HandleAuthErrors(ctx, err)
			return
		}
		ctx.JSON(
			http.StatusCreated,
			schemas.LoginResponse{Token: accessToken, RefreshToken: refreshToken},
		)
	}
}

// RefreshTokenHandler godoc
// @Summary Refresh Access Token
// @Description Refresh Access Token using a valid refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param token body schemas.RefreshTokenRequest true "Refresh Token"
// @Success 200 {object} schemas.LoginResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 401 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /v1/refresh [post]
func RefreshTokenHandler(authService *services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req schemas.RefreshTokenRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		accessToken, refreshToken, err := authService.RefreshAccessToken(req.RefreshToken)
		if err != nil {
			errors.HandleAuthErrors(ctx, err)
			return
		}

		ctx.JSON(
			http.StatusOK,
			schemas.LoginResponse{Token: accessToken, RefreshToken: refreshToken},
		)
	}
}

func AddAuthRoutes(router *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {
	authService := services.NewAuthService(db)

	router.POST("/login", LoginUserHandler(authService))
	router.POST("/refresh", RefreshTokenHandler(authService))
	return router
}
