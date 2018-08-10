package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"theAmazingCodeExample/app/common"
	"theAmazingCodeExample/app/models"
	"theAmazingCodeExample/app/security"
	"time"
)

func Signup(c *gin.Context) {

	emailValue := c.PostForm("email")
	password := c.PostForm("password")
	name := c.PostForm("name")
	lastName := c.PostForm("last_name")

	if name == "" || lastName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Name or last name not submitted"})
		return
	}

	if err := Validate(password, emailValue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": err.Error(), "detail": err.Error()})
		return
	}

	hash, _ := security.HashPassword(password)

	newUser := models.User{
		Email:    emailValue,
		Name:     name,
		LastName: lastName,
		Password: hash,
		RoleID:   models.USER,
	}

	if err := newUser.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": err.Error(), "detail": err.Error()})
		return
	}

	//Create email confirmation code
	recoveryCode := time.Now().Unix()
	stringCode := strconv.Itoa(int(recoveryCode))[len(strconv.Itoa(int(recoveryCode)))-5:]

	newEmailConfirmation := models.EmailConfirmation{
		UserID: newUser.ID,
		Code:   stringCode,
	}

	if err := newEmailConfirmation.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something unexpected happened", "detail": err.Error()})
		return
	}

	//Send email verification code
	//emailSubject := "Confirmación de email"
	//emailMessage := "Tu código de confirmación es: " + stringCode
	//if sendEmail := sendgrid.SendGenericIndividualEmail(emailSubject, emailMessage, newUser); sendEmail != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"description": sendEmail.Error()})
	//	return
	//}

	//Login information
	token, err := security.CreateToken(newUser.ID, newUser.Name, newUser.LastName, newUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something unexpected happened", "detail": err})
		return
	}

	permissionList, err := models.GetUserPermissions(newUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something unexpected happened", "detail": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"description": map[string]interface{}{"token": token, "email": newUser.Email, "name": newUser.Name, "lastName": newUser.LastName, "id": newUser.GUID, "permissions": permissionList, "profilePicture": newUser.ProfilePicture.Url}})

}

func SendRecoveryMail(c *gin.Context) {

	email := c.PostForm("email")

	if err := validateLiteralEmail(email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": err.Error(), "detail": err.Error()})
		return
	}

	userData, found, err := models.GetUserByEmail(email)
	if found == false {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Email not registered", "detail": ""})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something unexpected happened", "detail": err.Error()})
		return
	}

	recoveryCode := time.Now().Unix()
	stringCode := strconv.Itoa(int(recoveryCode))[len(strconv.Itoa(int(recoveryCode)))-4:]

	//Modify user recovery code
	userData.PasswordRecoveryCode = stringCode
	if err := userData.Modify(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something unexpected happened", "detail": err.Error()})
		return
	}

	//emailSubject := "Recupera tu contraseña"
	//emailMessage := "Tu codigo de recuperación de contraseña es: " + stringCode
	//if sendEmail := sendgrid.SendGenericIndividualEmail(emailSubject, emailMessage, userData); sendEmail != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"description": sendEmail.Error()})
	//	return
	//}

	c.JSON(http.StatusCreated, gin.H{"description": "Password recovery code sent"})

}

func ChangePasswordFromRecoveryCode(c *gin.Context) {

	email := c.PostForm("email")
	newPassword := c.PostForm("password")
	recoveryCode := c.PostForm("code")

	if err := validatePassword(newPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": err.Error(), "detail": err.Error()})
		return
	}

	if recoveryCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Password recovery code not submitted"})
		return
	}

	userData, err := validateRecoveryData(email, recoveryCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": err.Error(), "detail": err.Error()})
		return
	}

	hash, err := security.HashPassword(newPassword)

	//Modify user
	userData.Password = hash
	userData.PasswordRecoveryCode = ""

	if err := userData.Modify(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something unexpected happened", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "Password changed succesfully"})

}

func GetUsers(c *gin.Context) {

	idParam := c.Query("id")
	phoneParam := c.Query("phone")
	emailParam := c.Query("email")
	nameParam := c.Query("name")
	lastNameParam := c.Query("lastName")
	roleIdsParam := c.Query("role")
	limit := c.MustGet("limit").(int)
	offset := c.MustGet("offset").(int)
	column := c.MustGet("column").(string)
	order := c.MustGet("order").(string)

	//If role filter was informed, make a slice
	var roleList []string
	if roleIdsParam != "" {
		roleList = strings.Split(roleIdsParam, ",")
	}

	userList, quantity, err := models.GetUsers(limit, offset, idParam, phoneParam, emailParam, nameParam, lastNameParam, roleList, column, order)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Something unexpected happened when obtaining users", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": map[string]interface{}{"users": userList, "quantity": quantity}})

}

func ModifyUser(c *gin.Context) {

	userGUID := c.Param("id")
	name, wasInformedName := c.GetPostForm("name")
	lastName, wasInformedLastName := c.GetPostForm("last_name")
	roleId, wasInformedRole := c.GetPostForm("role")

	//Get user by GUID
	userData, found, err := models.GetUserByGuid(userGUID)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something unexpected happened when trying to obtaing the user", "detail": err.Error()})
		return
	}

	if wasInformedName == true {
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"description": "Name cannot be empty"})
			return
		}

		userData.Name = name
	}

	if wasInformedLastName == true {
		if lastName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"description": "Last name cannot be empty"})
			return
		}

		userData.LastName = lastName
	}

	if wasInformedRole == true {
		if roleId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"description": "Role cannot be empty"})
			return
		}

		roleIdVal, err := common.StringToUint(roleId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"description": "Invalid role ID"})
			return
		}

		roleData, found, err := models.GetRoleById(roleIdVal)
		if found == false {
			c.JSON(http.StatusBadRequest, gin.H{"description": "Role not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"description": "Something unexpected happened al obtener el rol", "detail": err.Error()})
			return
		}

		userData.RoleID = roleData.ID
	}

	//Modify user
	if err := userData.Modify(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something unexpected happened when trying to modify the user", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": userData})

}

func DisableUser(c *gin.Context) {

	userGUID := c.Param("id")

	//Get user by GUID
	userData, found, err := models.GetUserByGuid(userGUID)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something unexpected happened when trying to obtain the user", "detail": err.Error()})
		return
	}

	//Disable user
	if err := userData.Disable(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something unexpected happened when disabling user", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "User disabled"})

}

func EnableUser(c *gin.Context) {

	userGUID := c.Param("id")

	//Get user by GUID
	userData, found, err := models.GetUserByGuid(userGUID)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something unexpected happened when trying to obtain the user", "detail": err.Error()})
		return
	}

	//Enable user
	if err := userData.Enable(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something unexpected happened when enabling the user", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "User enabled"})

}