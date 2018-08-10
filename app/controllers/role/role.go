package role

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"theAmazingCodeExample/app/models"
)

func GetRoles(c *gin.Context) {

	roleList, quantity, err := models.GetRoles(c.MustGet("limit").(int), c.MustGet("offset").(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": "false", "description": "Something unexpected happened", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "true", "description": map[string]interface{}{"roles": roleList, "quantity": quantity}})

}
