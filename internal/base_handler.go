package internal

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TransactionalHandler(db *gorm.DB, handler func(c *gin.Context, tx *gorm.DB)) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = db.Transaction(func(tx *gorm.DB) error {
			handler(c, tx)
			if len(c.Errors) > 0 {
				return c.Errors.Last().Err
			}
			return nil
		})
	}
}
