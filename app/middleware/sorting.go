package middleware

import (
	"github.com/gin-gonic/gin"
)

func Sort() gin.HandlerFunc {
	return func(c *gin.Context) {
		defaultSort(c)
	}
}

func defaultSort(c *gin.Context) {
	var column string = ""
	var order string = ""

	columnParam := c.Query("Column")
	orderParam := c.Query("Order")

	if columnParam != "" && orderParam != "" {
		column = columnParam
		order = orderParam
	}

	c.Set("column", column)
	c.Set("order", order)

	c.Next()

}
