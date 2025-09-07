package private

import (
	"fleet-pulse-users-service/internal"
	"fleet-pulse-users-service/internal/errors"
	"fleet-pulse-users-service/internal/models"
	"fleet-pulse-users-service/internal/schemas"
	"fleet-pulse-users-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SearchUsersHandler Search users godoc
// @Summary Search users
// @Tags Users Internal
// @Accept json
// @Produce json
// @Param accept body schemas.SearchUsersPayload true "Search users payload"
// @Success 200 {object} []schemas.UserResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Router /v1/internal/users/search [post]
func SearchUsersHandler(userServiceConstructor func(db *gorm.DB) *services.UserService) func(c *gin.Context, tx *gorm.DB) {
	return func(c *gin.Context, tx *gorm.DB) {
		var req schemas.SearchUsersPayload
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userService := userServiceConstructor(tx)
		users, err := userService.SearchUsersByIds(req.UserIDs)
		if err != nil {
			errors.HandleUserErrors(c, err)
			c.Error(err)
			return
		}
		responses := internal.Map(users, func(u *models.User) schemas.UserResponse {
			return schemas.UserResponse{
				ID:        u.ID,
				Email:     u.Email,
				FirstName: u.FirstName,
				LastName:  u.LastName,
			}
		})
		c.JSON(http.StatusOK, responses)
	}
}

func AddInternalUserRoutes(router *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {
	userServiceConstructor := func(db *gorm.DB) *services.UserService {
		return services.NewUserService(db)
	}
	router.POST("/users/search",
		internal.TransactionalHandler(db, SearchUsersHandler(userServiceConstructor)),
	)
	return router
}
