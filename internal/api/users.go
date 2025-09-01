package api

import (
	"fleet-pulse-users-service/internal"
	"fleet-pulse-users-service/internal/errors"
	"fleet-pulse-users-service/internal/middlewares"
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
func RegisterUserHandler(userServiceConstructor func(db *gorm.DB) *services.UserService) func(c *gin.Context, tx *gorm.DB) {
	return func(c *gin.Context, tx *gorm.DB) {
		var req schemas.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userService := userServiceConstructor(tx)
		createdUser, err := userService.RegisterNewUser(req)
		if err != nil {
			errors.HandleUserErrors(c, err)
			c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, schemas.UserResponse{
			ID:        createdUser.ID,
			Email:     createdUser.Email,
			FirstName: createdUser.FirstName,
			LastName:  createdUser.LastName,
		})
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
func GetCurrentUserHandler(userServiceConstructor func(db *gorm.DB) *services.UserService) func(c *gin.Context, tx *gorm.DB) {
	return func(c *gin.Context, tx *gorm.DB) {
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

		userService := userServiceConstructor(tx)
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

// AcceptInviteHandler godoc
// @Summary Accept user invite
// @Description Accept an invite and set password
// @Tags Users
// @Accept json
// @Produce json
// @Param accept body schemas.AcceptInviteRequest true "Accept invite request"
// @Success 200 {object} schemas.UserResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /v1/users/invite/accept [post]
func AcceptInviteHandler(userServiceConstructor func(db *gorm.DB) *services.UserService) func(c *gin.Context, tx *gorm.DB) {
	return func(c *gin.Context, tx *gorm.DB) {
		var req schemas.AcceptInviteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userService := userServiceConstructor(tx)
		user, err := userService.AcceptInvite(req.Token, req.Password)
		if err != nil {
			errors.HandleUserErrors(c, err)
			return
		}

		c.JSON(http.StatusOK, schemas.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		})
	}
}

func AddUserRoutes(router *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {
	userServiceConstructor := func(db *gorm.DB) *services.UserService {
		return services.NewUserService(db)
	}

	router.POST("/users",
		internal.TransactionalHandler(db, RegisterUserHandler(userServiceConstructor)),
	)

	router.GET("/users/current",
		middlewares.JWTAuthMiddleware(services.NewAuthService(db)),
		internal.TransactionalHandler(db, GetCurrentUserHandler(userServiceConstructor)),
	)

	router.POST("/users/invite/accept",
		internal.TransactionalHandler(db, AcceptInviteHandler(userServiceConstructor)),
	)

	return router
}
