package main

import (
	"net/http"
	"restapi/controllers"

	"github.com/gin-gonic/gin"
	validator "github.com/gobeam/custom-validator"
)

func main() {
	router := gin.Default()

	validate := []validator.ExtraValidation{
		{Tag: "number", Message: "Invalid %s Format!"},
		{Tag: "email", Message: "Invalid %s Format!"},
	}

	validator.MakeExtraValidation(validate)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello jay",
		})
	})

	// users
	router.Use(validator.Errors())
	{
		router.POST("/register", controllers.UserRegister)
		router.POST("/register/guide", controllers.GuideRegister)
		router.POST("/login", controllers.UserLogin)
	}
	router.GET("/users", controllers.Users)
	router.GET("/confirm-email/:token", controllers.ConfirmEmail)
	router.POST("/resend-email", controllers.ResendEmail)

	router.GET("/create-cookie", controllers.Cookie)
	router.GET("/delete-cookie", controllers.DeleteCookie)

	router.Run()
}
