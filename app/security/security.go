package security

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"theAmazingCodeExample/app/config"
	"theAmazingCodeExample/app/models"
	"time"
)

func Login(c *gin.Context) {

	email := c.PostForm("email")
	password := c.PostForm("password")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"description": "No email submitted", "detail": ""})
		return
	}
	if password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"description": "No password submitted", "detail": ""})
		return
	}

	u, found, err := models.GetUserByEmail(email)
	if found == false {
		c.JSON(http.StatusUnauthorized, gin.H{"description": "Email not registered", "detail": ""})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something unexpected happened", "detail": err.Error()})
		return
	}

	if CheckPasswordHash(password, u.Password) == false {
		c.JSON(http.StatusUnauthorized, gin.H{"description": "Invalid email or password", "detail": ""})
		return
	}

	if u.Enabled == false{
		c.JSON(http.StatusForbidden, gin.H{"description": "Your user was disabled by an administrator", "detail": ""})
		return
	}

	token, err := CreateToken(u.ID, u.Name, u.LastName, u.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something unexpected happenedo", "detail": err.Error()})
		return
	}

	permissionList, err := models.GetUserPermissions(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something unexpected happened", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": "true", "description": map[string]interface{}{"token": token, "email": u.Email, "name": u.Name, "lastName": u.LastName, "id": u.GUID, "permissions": permissionList, "profilePicture": u.ProfilePicture.Url}})

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type JWTCustomClaims struct {
	Id       uint
	Email    string
	Name     string
	LastName string
	jwt.StandardClaims
}

type JWTToken struct {
	Id       uint
	Name     string
	LastName string
	Email    string
}

func CreateToken(id uint, name string, lastName string, email string) (string, error) {
	claims := JWTCustomClaims{
		id,
		email,
		name,
		lastName,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			Issuer:    name + " " + lastName,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.GetConfig().JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetTokenData(tokenString string) (JWTToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JWT_SECRET), nil
	})
	if err != nil {
		return JWTToken{}, err
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return JWTToken{
			Id:       claims.Id,
			Name:     claims.Name,
			LastName: claims.LastName,
			Email:    claims.Email,
		}, nil
	} else {
		return JWTToken{}, errors.New("Invalid token")
	}
}
