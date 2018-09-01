package user

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"theAmazingCodeExample/app/common"
	"theAmazingCodeExample/app/helpers/amazonS3"
	"theAmazingCodeExample/app/models"
	"theAmazingCodeExample/app/security"
	"theAmazingCodeExample/app/services/theAmazingEmailSender"
	"theAmazingCodeExample/app/services/theAmazingSmsSender"
	"time"
)

func SendConfirmationEmail(c *gin.Context) {

	userID := c.MustGet("id").(uint)

	//Get user data
	userData, found, err := models.GetUserById(userID)
	if found == false {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Get user email confirmations (if he has any)
	userEmailConfirmation, found, err := userData.GetEmailConfirmation()
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Your user doesn't have a pending email change"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Create email confirmation code
	recoveryCode := time.Now().Unix()
	stringCode := strconv.Itoa(int(recoveryCode))[len(strconv.Itoa(int(recoveryCode)))-5:]

	//Send email verification code
	emailSubject := "Confirmación de email"
	emailMessage := "Tu código de confirmación es: " + stringCode
	if sendEmail := theAmazingEmailSender.SendGenericIndividualEmail(emailSubject, emailMessage, userData); sendEmail != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": sendEmail.Error()})
		return
	}

	//Modify email confirmation
	userEmailConfirmation.Code = stringCode
	if err := userEmailConfirmation.Modify(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "Confirmation email sent"})

}

func VerifyEmail(c *gin.Context) {

	userID := c.MustGet("id").(uint)
	emailCode := c.PostForm("code")

	//Get user data
	userData, found, err := models.GetUserById(userID)
	if found == false {
		c.JSON(http.StatusNotFound, gin.H{"description": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Check email code
	if emailCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Email confirmation code not sent"})
		return
	}

	//Get user email confirmations (if he has any)
	userEmailConfirmation, found, err := userData.GetEmailConfirmation()
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "You don't have a pending email change"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Check confirmation code matches
	if userEmailConfirmation.Code != emailCode {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Invalid change code"})
		return
	}

	//Delete email confirmations of the user
	if err := userEmailConfirmation.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": err.Error()})
		return
	}

	//Modify user
	userData.Email = userEmailConfirmation.Email
	if err := userData.Modify(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "Email changed successfully"})

}

func ModifyEmail(c *gin.Context) {

	userID := c.MustGet("id").(uint)
	email := c.PostForm("email")

	//Validate Email
	if err := validateLiteralEmail(email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": err.Error(), "detail": err.Error()})
		return
	}

	//Check email is not being used
	isBeingUsed, err := models.CheckEmailExistence(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}
	if isBeingUsed == true {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Email already in use"})
		return
	}

	//Get user data
	userData, found, err := models.GetUserById(userID)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//If email changes, register the change and send confirmation email
	if userData.Email != email {

		//Delete previous email confirmations
		if err := userData.DeleteEmailConfirmations(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"description": err.Error()})
			return
		}

		//Create email confirmation
		recoveryCode := time.Now().Unix()
		stringCode := strconv.Itoa(int(recoveryCode))[len(strconv.Itoa(int(recoveryCode)))-5:]

		newEmailConfirmation := models.EmailConfirmation{
			UserID: userData.ID,
			Code:   stringCode,
			Email:  email,
		}

		if err = newEmailConfirmation.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
			return
		}

		//Send email verification code
		emailSubject := "Confirmación de email"
		emailMessage := "Tu código de confirmación es: " + stringCode
		if sendEmail := theAmazingEmailSender.SendGenericIndividualEmail(emailSubject, emailMessage, userData); sendEmail != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"description": sendEmail.Error()})
			return
		}

	} else {

		c.JSON(http.StatusBadRequest, gin.H{"description": "You have submitted the same email that you have"})
		return

	}

	c.JSON(http.StatusOK, gin.H{"description": "Email confirmation mail sent"})

}

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
		Email:  emailValue,
	}

	if err := newEmailConfirmation.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Sendgrid function needs this value
	newUser.Email = emailValue

	//Send email verification code
	emailSubject := "Confirmación de email"
	emailMessage := "Ingresa al siguiente link para confirmar tu contraseña: http://localhost:5000/confirmEmail?code=" + stringCode + "&email=" + emailValue
	if sendEmail := theAmazingEmailSender.SendGenericIndividualEmail(emailSubject, emailMessage, newUser); sendEmail != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": sendEmail.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "Account created. Please confirm your email before you can access the platform"})
}

func SendRecoveryMail(c *gin.Context) {

	email := c.PostForm("email")

	if err := validateLiteralEmail(email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": err.Error(), "detail": err.Error()})
		return
	}

	userData, found, err := models.GetUserByEmail(email)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Email not registered", "detail": ""})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	recoveryCode := time.Now().Unix()
	stringCode := strconv.Itoa(int(recoveryCode))[len(strconv.Itoa(int(recoveryCode)))-4:]

	//Modify user recovery code
	userData.PasswordRecoveryCode = stringCode
	if err := userData.Modify(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	emailSubject := "Recupera tu contraseña"
	emailMessage := "Tu codigo de recuperación de contraseña es: " + stringCode
	if sendEmail := theAmazingEmailSender.SendGenericIndividualEmail(emailSubject, emailMessage, userData); sendEmail != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": sendEmail.Error()})
		return
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"description": "Something went wrong when obtaining users", "detail": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong when trying to obtaing the user", "detail": err.Error()})
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
			c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong when obtaining the role", "detail": err.Error()})
			return
		}

		userData.Role = roleData
	}

	//Modify user
	if err := userData.Modify(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong when trying to modify the user", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": userData})

}

func EnableUser(c *gin.Context) {

	userGUID := c.Param("id")
	enable := c.PostForm("enabled")

	//Get user by GUID
	userData, found, err := models.GetUserByGuid(userGUID)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong when trying to obtain the user", "detail": err.Error()})
		return
	}

	enabledValue, err := strconv.Atoi(enable)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	if enabledValue == 1 {
		userData.Enabled = true
	} else {
		userData.Enabled = false
	}

	if err := userData.Modify(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong when enabling the user", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": userData})

}

func ChangePassword(c *gin.Context) {

	userID := c.MustGet("id").(uint)

	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")

	userData, found, err := models.GetUserById(userID)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Check old password matches the actual password
	if match := security.CheckPasswordHash(oldPassword, userData.Password); !match {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Invalid password"})
		return
	}

	//Validate new password is valid
	if err := validatePassword(newPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": err.Error(), "detail": err.Error()})
		return
	}

	newHash, _ := security.HashPassword(newPassword)

	userData.Password = newHash

	if err := userData.Modify(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "Password changed successfully"})

}

func ModifyUserName(c *gin.Context) {

	userID := c.MustGet("id").(uint)

	firstName := c.PostForm("name")
	lastName := c.PostForm("last_name")

	if firstName == "" || lastName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Name or last name were not submitted"})
		return
	}

	//Get user data
	userData, found, err := models.GetUserById(userID)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	userData.Name = firstName
	userData.LastName = lastName

	if err := userData.Modify(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": userData})

}

func GetUserProfile(c *gin.Context) {

	userID := c.MustGet("id").(uint)

	userData, found, err := models.GetUserById(userID)
	if found == false {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": userData})

}

func AddProfilePicture(c *gin.Context) {

	userID := c.MustGet("id").(uint)
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	photoFile := form.File["profile_picture"][0]

	//Get user data
	userData, found, err := models.GetUserById(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "User not found"})
		return
	}

	//Check if user has a profile picture. If so, delete old picture
	if userData.ProfilePicture.ID != 0 {
		if err = amazonS3.DeletePictureFromS3(userData.ProfilePicture); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
			return
		}
	}

	tasks := []*amazonS3.UploadImageTask{
		{
			FileHeader: photoFile,
			UserID:     userID,
			Function:   uploadAndSaveProfilePicture,
		},
	}

	p := amazonS3.NewPool(tasks, 3)
	p.Run()

	var numErrors = 0
	for _, task := range p.Tasks {
		if task.Err != nil {
			println(task.Err.Error())
			numErrors++
		}
	}

	if numErrors > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Unexpected error when saving the profile picture"})
	} else {
		c.JSON(http.StatusOK, gin.H{"description": "Profile picture added"})
	}

}

func DeleteProfilePicture(c *gin.Context) {

	userID := c.MustGet("id").(uint)

	//Get user data
	userData, found, err := models.GetUserById(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "User not found"})
		return
	}

	//Check if user has a profile picture
	if userData.ProfilePicture.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"description": "You don't have a profile picture"})
		return
	}

	if err = amazonS3.DeletePictureFromS3(userData.ProfilePicture); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "Profile picture deleted"})

}

func uploadAndSaveProfilePicture(header *multipart.FileHeader, userID uint) error {

	s3key, url, err := amazonS3.UploadImageToS3(header)
	if err != nil {
		return err
	}

	newProfilePicture := models.ProfilePicture{
		Url:    url,
		S3Key:  s3key,
		UserID: userID,
	}

	if err = newProfilePicture.Save(); err != nil {
		return err
	}

	return nil

}

func SendVerificationSMS(c *gin.Context) {

	userID := c.MustGet("id").(uint)

	//Get user data
	userData, found, err := models.GetUserById(userID)
	if found == false {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Get user phone confirmations (if he has any)
	userPhoneConfirmation, found, err := userData.GetPhoneConfirmation()
	if found == false {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Your user doesn't have any pending phone confirmations"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Create phone confirmation code
	recoveryCode := time.Now().Unix()
	phoneCode := strconv.Itoa(int(recoveryCode))[len(strconv.Itoa(int(recoveryCode)))-4:]

	//Create task to send verification code
	if err := theAmazingSmsSender.SendSms(userPhoneConfirmation.Phone, phoneCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Modify phone confirmation
	userPhoneConfirmation.Code = phoneCode
	if err := userPhoneConfirmation.Modify(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "SMS sent"})

}

func ModifyPhone(c *gin.Context) {

	userID := c.MustGet("id").(uint)
	phoneNumber := c.PostForm("phone_number")

	//Check fields
	if phoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"description": "No phone submitted"})
		return
	}

	//Check phone number is valid
	isValidPhoneNumber, err, phoneData := theAmazingSmsSender.ValidatePhoneNumber(phoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}
	if isValidPhoneNumber == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Invalid phone number"})
		return
	}
	if phoneData.CountryCode != "AR" {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Only argentinian phone numbers are available"})
		return
	}

	//Check the cellphone is not being used by anyone else
	phoneUsage, err := models.CheckPhoneUsage(phoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}
	if phoneUsage == true {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Phone number already in use"})
		return
	}

	//Get user
	userData, found, err := models.GetUserById(userID)
	if found == false {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Delete previous phone confirmations
	if err := userData.DeletePhoneConfirmations(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": err.Error()})
		return
	}

	//Create phone confirmation
	randomCode := time.Now().Unix()
	phoneCode := strconv.Itoa(int(randomCode))[len(strconv.Itoa(int(randomCode)))-4:]

	newPhoneConfirmation := models.PhoneConfirmation{
		UserID: userData.ID,
		Phone:  phoneNumber,
		Code:   phoneCode,
	}

	if err = newPhoneConfirmation.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Send verification code
	if err := theAmazingSmsSender.SendSms(phoneNumber, phoneCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "A SMS code has been sent to your phone"})

}

func ConfirmPhoneCode(c *gin.Context) {

	userID := c.MustGet("id").(uint)
	confirmationCode := c.PostForm("code")

	//Check fields
	if confirmationCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Confirmation code hasn't been submitted"})
		return
	}

	//Get user
	userData, found, err := models.GetUserById(userID)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Get user phone confirmations (if he has any)
	userPhoneConfirmation, found, err := userData.GetPhoneConfirmation()
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Your user doesn't have any pending phone change"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Check confirmation code matches
	if userPhoneConfirmation.Code != confirmationCode {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Invalid confirmation code"})
		return
	}

	//Delete phone confirmations of the user
	if err := userPhoneConfirmation.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": err.Error()})
		return
	}

	//Modify user
	userData.Phone = userPhoneConfirmation.Phone
	if err := userData.Modify(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "Phone modified successfully"})

}

func ConfirmEmail(c *gin.Context) {

	email := c.Query("email")
	code := c.Query("code")

	if email == "" || code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Email or code missing"})
		return
	}

	//Get user email confirmations by code
	userEmailConfirmation, found, err := models.GetEmailConfirmationByCode(code)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Invalid code"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Get user data
	userData, found, err := models.GetUserById(userEmailConfirmation.UserID)
	if found == false {
		c.JSON(http.StatusNotFound, gin.H{"description": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Modify user
	userData.Email = userEmailConfirmation.Email
	if err := userData.Modify(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Delete email confirmations of the user
	if err := userEmailConfirmation.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "Email confirmed successfully. You can now log in"})

}
