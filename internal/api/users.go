package api

import (
	"fleet-pulse-users-service/internal/errors"
	"fleet-pulse-users-service/internal/middlewares"
	"fleet-pulse-users-service/internal/repositories"
	"fleet-pulse-users-service/internal/schemas"
	"fleet-pulse-users-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RegisterUserHandler godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags Users
// @Accept json
// @Produce json
// @Param user body schemas.CreateUserRequest true "User data"
// @Success 201 {object} schemas.UserResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 409 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /v1/users [post]
func RegisterUserHandler(userService *services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req schemas.CreateUserRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createdUser, err := userService.RegisterNewUser(req)
		if err != nil {
			errors.HandleUserErrors(ctx, err)
			return
		}
		ctx.JSON(
			http.StatusCreated,
			schemas.UserResponse{
				ID:        createdUser.ID,
				Email:     createdUser.Email,
				FirstName: createdUser.FirstName,
				LastName:  createdUser.LastName,
			},
		)
	}
}

// GetCurrentUserHandler godoc
// @Summary Get Current User
// @Description Get information about the currently authenticated user
// @Tags Users
// @Produce json
// @Success 200 {object} schemas.UserResponse
// @Failure 401 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Router /v1/users/current [get]
// @Security Bearer
func GetCurrentUserHandler(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("current_user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userUUID, err := uuid.Parse(userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		user, err := userService.GetUserById(userUUID)
		if err != nil || user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, schemas.UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		})
	}
}

func AddUserRoutes(router *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {
	userRepo := repositories.NewUserRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(refreshTokenRepo, userRepo)

	router.POST("/users", RegisterUserHandler(userService))
	router.GET("/users/current",
		middlewares.JWTAuthMiddleware(authService),
		GetCurrentUserHandler(userService),
	)
	return router
}
