package schemas

type RegisterUser struct {
	Username        string `json:"username" form:"username" binding:"required,min=3,max=100"`
	Email           string `json:"email" form:"email" binding:"required,email"`
	Phone           string `json:"phone" form:"phone" binding:"required,number,min=10,max=20"`
	Password        string `json:"password" form:"password" binding:"required,min=6,max=100"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required,eqfield=Password"`
}

type RegisterGuide struct {
	Register RegisterUser
	Image    string `json:"image" form:"image"`
	Region   string `json:"region" form:"region" binding:"required"`
	Role     string `json:"role" form:"role" binding:"required"`
	RegisterTranslator
}

type RegisterTranslator struct {
	Language string `json:"language" form:"language" binding:"required"`
}

type LoginUser struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=6,max=100"`
}
