package controllers

import (
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"restapi/config"
	"restapi/libs"
	"restapi/models"
	"restapi/schemas"
	"runtime"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 4 // 4 MB

type bodylink struct {
	Name string
	URL  string
}

func UserRegister(c *gin.Context) {
	var payload schemas.RegisterUser

	// validation
	if err := c.ShouldBindWith(&payload, binding.Form); err != nil {
		_ = c.AbortWithError(422, err).SetType(gin.ErrorTypeBind)
		return
	}

	// check duplicate email
	email := models.CheckEmailExists(payload.Email)
	if email != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email is already taken."})
		return
	}

	phone := models.CheckPhoneExists(payload.Phone)
	if phone != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Phone number is already taken."})
		return
	}

	// insert user
	u, userErr := models.SaveUser(payload)

	if userErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed creating account"})
		return
	}

	// send email after creating account
	token := uuid.New().String()
	_ = models.SaveConfirmation(token, u.ID)

	link := os.Getenv("APP_URL") + "/confirm-email/" + token
	templateData := bodylink{
		Name: u.Username,
		URL:  link,
	}

	runtime.GOMAXPROCS(1)
	go libs.SendEmailVerification(payload.Email, templateData)

	c.JSON(http.StatusCreated, gin.H{"message": "Success, check your email to verification."})
}

func GuideRegister(c *gin.Context) {
	var payload schemas.RegisterGuide

	// validation
	if err := c.ShouldBindWith(&payload, binding.Form); err != nil {
		_ = c.AbortWithError(422, err).SetType(gin.ErrorTypeBind)
		return
	}

	// check duplicate email
	email := models.CheckEmailExists(payload.Register.Email)
	if email != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email is already taken."})
		return
	}

	phone := models.CheckPhoneExists(payload.Register.Phone)
	if phone != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Phone is already taken."})
		return
	}

	// form image
	imgExt := []string{"JPEG", "PNG"}

	form, _ := c.MultipartForm()

	licenses := form.File["image"]
	if len(licenses) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Image is required."})
		return
	}

	var files []*multipart.FileHeader
	for _, file := range licenses {
		ext := filepath.Ext(file.Filename)

		// check image format
		imgFormat, _ := imaging.FormatFromExtension(ext)
		checkImg := CheckImage(imgExt, imgFormat.String())

		// check image size
		if file.Size > MAX_UPLOAD_SIZE {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Image must less than 4 MB."})
			return
		}

		if checkImg == false {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Image must be jpeg, jpg or png."})
			return
		}

		files = append(files, file)
	}

	// save image
	var fileName []string
	for _, img := range files {
		ext := filepath.Ext(img.Filename)

		imgName := uuid.New().String() + ext
		dst := "./static/licenses/" + imgName
		c.SaveUploadedFile(img, dst)

		fileName = append(fileName, imgName)
	}
	payload.Image = strings.Join(fileName, ",")

	// insert user
	u, userErr := models.SaveDriver(payload)

	if userErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed creating account"})
		return
	}

	// send email after creating account
	token := uuid.New().String()
	_ = models.SaveConfirmation(token, u.ID)

	link := os.Getenv("APP_URL") + "/confirm-email/" + token
	templateData := bodylink{
		Name: u.Username,
		URL:  link,
	}

	runtime.GOMAXPROCS(1)
	go libs.SendEmailVerification(payload.Register.Email, templateData)

	c.JSON(http.StatusCreated, gin.H{"message": "Success, check your email to verification."})

}

func CheckImage(format []string, extension string) bool {
	for _, img := range format {
		if img == extension {
			return true
		}
	}
	return false
}

func UserLogin(c *gin.Context) {
	var payload schemas.LoginUser

	// validation
	if err := c.ShouldBindWith(&payload, binding.JSON); err != nil {
		_ = c.AbortWithError(422, err).SetType(gin.ErrorTypeBind)
		return
	}

	user, errLogin := models.VerifyLogin(payload.Email, payload.Password)

	if errLogin != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials."})
		return
	}

	confirmation, _ := models.FindConfirmation(user.ID)

	// create token
	token, _ := config.CreateToken(user.ID)

	tokens := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}

	if !confirmation.Activated {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Account is not actived, check your email to verified"})
	} else {
		c.JSON(http.StatusOK, tokens)
	}

}

func ConfirmEmail(c *gin.Context) {
	id := c.Param("token")

	confirmation := models.CheckActivatedById(id)
	if confirmation.Activated == true {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Your account already activated."})
		return
	}

	_ = models.ConfirmEmailModel(id)

	// create token
	token, _ := config.CreateToken(confirmation.User_id)

	DOMAIN := "localhost"
	// add jwt to cookies
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", token.RefreshToken, 900, "/", DOMAIN, true, true)
	c.SetCookie("refresh_token", token.RefreshToken, 900, "/", DOMAIN, true, true)

	// c.Redirect(http.StatusPermanentRedirect, "/")

	c.JSON(http.StatusOK, gin.H{"message": "Account verified,log in."})
}

func Cookie(c *gin.Context) {
	// c.SetSameSite(http.SameSiteLaxMode)
	cookie := &http.Cookie{}
	cookie.Name = "tes"
	cookie.Value = "ini cookie"
	cookie.MaxAge = 900
	cookie.Path = "/create-cookie"
	cookie.Domain = "localhost"
	cookie.Secure = true
	cookie.HttpOnly = false

	http.SetCookie(c.Writer, cookie)

	co, _ := c.Cookie("tes")
	log.Println("cookie --->", co)

}

func DeleteCookie(c *gin.Context) {
	c.SetCookie("tes", "ini cookie", -1, "/create-cookie", "localhost", true, false)
}

func ResendIsExpired(resend_expired int64) bool {

	return int64(time.Now().Unix()) > resend_expired
}

func ResendEmail(c *gin.Context) {
	var payload schemas.ResendEmail

	// validation
	if err := c.ShouldBindWith(&payload, binding.JSON); err != nil {
		_ = c.AbortWithError(422, err).SetType(gin.ErrorTypeBind)
		return
	}

	// check email exist in database
	email := models.CheckEmailExists(payload.Email)
	if email != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email not found."})
		return
	}

	// check email confirmation activated or not
	user := models.FindUserByEmail(payload.Email)

	confirmation, _ := models.FindConfirmation(user.ID)

	if confirmation.Activated == true {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Your account already activated."})
		return
	}

	if confirmation.Resend_expired == 0 || ResendIsExpired(confirmation.Resend_expired) {
		// send email after creating account
		token := uuid.New().String()
		// add 5 minute resend expired
		models.SaveResendEmail(token, user.ID)

		link := os.Getenv("APP_URL") + "/confirm-email/" + token
		templateData := bodylink{
			Name: user.Username,
			URL:  link,
		}

		runtime.GOMAXPROCS(1)
		go libs.SendEmailVerification(payload.Email, templateData)

		c.JSON(http.StatusOK, gin.H{"message": "Email confirmation has send"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You can try 5 minute later."})
	}

}

func Users(c *gin.Context) {
	users := models.GetUsers()
	c.JSON(http.StatusOK, gin.H{"users": users})
}
