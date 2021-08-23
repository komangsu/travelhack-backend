package schemas

import (
	"time"
)

type UserSchema struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username" binding:"required,min=3,max=100"`
	Email     string    `json:"email" binding:"required,email"`
	Password  string    `json:"password" binding:"required,min=6,max=100"`
	Phone     string    `json:"phone" binding:"required,number,min=10,max=10"`
	Avatar    string    `json:"avatar"`
	Image     string    `json:"image,omitempty"`
	Role      string    `json:"role"`
	Region    string    `json:"region,omitempty" binding:"required"`
	Language  string    `json:"language,omitempty" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResendEmail struct {
	Email string `json:"email" binding:"required,email"`
}
