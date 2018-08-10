package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

func Paginate() gin.HandlerFunc {
	return func(c *gin.Context) {
		defaultPaginate(c)
	}
}

func defaultPaginate(c *gin.Context) {
	var offset int = -1
	var limit int = -1

	limitParam := c.Query("Limit")
	offsetParam := c.Query("Offset")

	if limitParam != "" {
		tmp, err := strconv.ParseInt(limitParam, 10, 32)
		limit = int(tmp)
		if err != nil || !genIsLimitValid(limit) {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid limit"})
			c.Abort()
			return
		}
	}
	c.Set("limit", limit)

	if limit != -1 && offsetParam != "" {
		tmp, err := strconv.ParseInt(offsetParam, 10, 32)
		offset = int(tmp)
		if err != nil || !genIsOffsetValid(offset) {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid offset"})
			c.Abort()
			return
		}
	}
	c.Set("offset", offset)

	c.Next()
}

func genIsOffsetValid(o int) bool {
	if o >= 0 {
		return true
	}
	return false
}

func genIsLimitValid(l int) bool {
	if l > 0 {
		return true
	}
	return false
}

func ApplyPaginate(c *gin.Context, db *gorm.DB) *gorm.DB {

	offset := c.MustGet("offset").(int)
	limit := c.MustGet("limit").(int)

	db = db.Offset(offset).Limit(limit)

	return db
}
