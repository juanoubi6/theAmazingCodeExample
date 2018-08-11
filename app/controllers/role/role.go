package role

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"theAmazingCodeExample/app/models"
	"theAmazingCodeExample/app/common"
)

func GetRoles(c *gin.Context) {

	roleList, quantity, err := models.GetRoles(c.MustGet("limit").(int), c.MustGet("offset").(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": "false", "description": "Something unexpected happened", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "true", "description": map[string]interface{}{"roles": roleList, "quantity": quantity}})

}

func ModifyPermissions(c *gin.Context){

	id := c.Param("id")
	permissions,_ := c.GetPostFormArray("permission")

	idVal,err := common.StringToUint(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Invalid role id", "detail": err.Error()})
		return
	}

	//Get role by ID
	roleData, found, err := models.GetRoleById(idVal)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Role not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something unexpected happened when trying to obtain the user", "detail": err.Error()})
		return
	}

	if err = roleData.ReplacePermissions(permissions); err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something unexpected happened when trying to modify role permissions", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"description": roleData})
	return

}
