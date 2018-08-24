package oauthGoogle

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"theAmazingCodeExample/app/config"
	"theAmazingCodeExample/app/models"
	"theAmazingCodeExample/app/security"
)

var oauthStateString = "iwillberandomsomeday"

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:5000/googleCallback",
	ClientID:     config.GetConfig().GOOGLE_CLIENT_ID,
	ClientSecret: config.GetConfig().GOOGLE_CLIENT_SECRET,
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

type GoogleData struct {
	Name     string `json:"family_name"`
	LastName string `json:"given_name"`
	Email    string `json:"email"`
	GoogleID string `json:"id"`
	Picture  string `json:"picture"`
}

func MainPage(c *gin.Context) {
	var htmlIndex = `<html><body><a href="/login">Google Log In</a></body></html>`
	c.Header("Content-Type", "text/html")
	c.String(200, "<html>%s</html>", htmlIndex)
}

func RedirectToGoogle(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleGoogleCallback(c *gin.Context) {
	content, err := getUserInfo(c.Query("state"), c.Query("code"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	var googleData GoogleData
	err = json.Unmarshal(content, &googleData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	var userData models.User

	//Check if user already exists. If not, create it.
	userData, found, err := models.GetUserByGoogleId(googleData.GoogleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}
	if found == false {

		//Check if there is an user with this email
		userData, found, err = models.GetUserByEmail(googleData.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
			return
		}
		if found == true {
			userData.GoogleID = googleData.GoogleID
			if err := userData.Modify(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"description": err.Error(), "detail": err.Error()})
				return
			}
		} else {
			userData = models.User{
				Email:    googleData.Email,
				Name:     googleData.Name,
				LastName: googleData.LastName,
				RoleID:   models.USER,
				GoogleID: googleData.GoogleID,
			}

			if err := userData.Save(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"description": err.Error(), "detail": err.Error()})
				return
			}
		}
	}

	//Login information
	token, err := security.CreateToken(userData.ID, userData.Name, userData.LastName, userData.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err})
		return
	}

	permissionList, err := models.GetUserPermissions(userData.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"description": map[string]interface{}{"token": token, "email": userData.Email, "name": userData.Name, "lastName": userData.LastName, "id": userData.GUID, "permissions": permissionList, "profilePicture": userData.ProfilePicture.Url}})
	return
}

func getUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, errors.New("Invalid oauth state")
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, errors.New("Code exchange failed: " + err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, errors.New("Failed getting user info: " + err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("Failed reading response body: " + err.Error())
	}

	return contents, nil
}
