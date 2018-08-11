package router

import (
	"github.com/aviddiviner/gin-limit"
	"github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"theAmazingCodeExample/app/config"
	"theAmazingCodeExample/app/controllers/address"
	"theAmazingCodeExample/app/controllers/role"
	"theAmazingCodeExample/app/controllers/user"
	"theAmazingCodeExample/app/middleware"
	"theAmazingCodeExample/app/migrations"
	"theAmazingCodeExample/app/security"
)

var router *gin.Engine

func CreateRouter() {
	router = gin.New()

	router.Use(gin.Logger())
	router.Use(nice.Recovery(recoveryHandler))
	router.Use(limit.MaxAllowed(10))
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET,PUT,POST,DELETE"},
		AllowHeaders:    []string{"accept,x-access-token,content-type,authorization"},
	}))

	public := router.Group("/")
	{
		public.GET("/migrations", migrations.Run)
	}

	/* Routes */
	loginRoutes := router.Group("/")
	{
		loginRoutes.POST("/login", security.Login)
		loginRoutes.POST("/signup", user.Signup)
		loginRoutes.POST("/recoverPassword", user.SendRecoveryMail)
		loginRoutes.PUT("/password", user.ChangePasswordFromRecoveryCode)
	}

	userManagment := router.Group("/users", middleware.IsAdmin(), middleware.ValidateTokenAndPermission("User Management"))
	{
		userManagment.GET("", middleware.Paginate(), middleware.Sort(), user.GetUsers)
		userManagment.PUT("/:id", user.ModifyUser)
		userManagment.PUT("/:id/enable", user.EnableUser)
	}

	rolesManagment := router.Group("/roles", middleware.IsAdmin(), middleware.ValidateTokenAndPermission("Role Management"))
	{
		rolesManagment.GET("", middleware.Paginate(), role.GetRoles)
		rolesManagment.PUT("/:id/permissions", role.ModifyPermissions)
	}

	userProfile := router.Group("/user", middleware.ValidateTokenAndPermission("Profile"))
	{
		//User profile endpoints
		userProfile.PUT("/profile", user.ModifyUserName)
		userProfile.GET("/profile", user.GetUserProfile)
		userProfile.POST("/profile/picture", user.AddProfilePicture)
		userProfile.DELETE("/profile/picture", user.DeleteProfilePicture)
		userProfile.PUT("/password", user.ChangePassword)

		//Address endpoints
		userProfile.GET("/address", address.GetAddresses)
		userProfile.POST("/address", address.AddAddress)
		userProfile.PUT("/address/:id", address.ModifyAddress)
		userProfile.PUT("/address/:id/main", address.MarkAsMain)
		userProfile.DELETE("/address/:id", address.DeleteAddress)

		//Phone endpoints
		//userProfile.POST("/phone/verificationSMS", user.ModifyPhone)
		//userProfile.POST("/phone", user.ConfirmPhoneCode)
		//userProfile.GET("/resendVerificationSMS", user.SendVerificationSMS)

		//Email change and verification
		userProfile.PUT("/profile/email", user.ModifyEmail)
		userProfile.PUT("/verifyEmail", user.VerifyEmail)
		userProfile.GET("/resendConfirmationEmail", user.SendConfirmationEmail)

	}

}

func RunRouter() {
	router.Run(":" + config.GetConfig().PORT)
}

func recoveryHandler(c *gin.Context, err interface{}) {
	detail := ""
	if config.GetConfig().ENV == "develop" {
		detail = err.(error).Error()
	}
	c.JSON(http.StatusInternalServerError, gin.H{"success": "false", "description": detail})
}
