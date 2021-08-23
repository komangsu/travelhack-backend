package models

import (
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"restapi/config"
	"restapi/schemas"
	"time"
)

type ListUser struct {
	Id         uint64    `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Avatar     string    `json:"avatar"`
	Image      string    `json:"image,omitempty"`
	Region     string    `json:"region,omitempty"`
	Language   string    `json:"language,omitempty"`
	Role       string    `json:"role"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}

type Claims struct {
	UserId string `json:"email"`
	jwt.StandardClaims
}

var (
	JWT_SECRET_KEY = []byte(os.Getenv("JWT_SECRET_KEY"))
)

// hash password
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// User guest
func SaveUser(payload schemas.RegisterUser) (schemas.UserSchema, error) {
	var u schemas.UserSchema

	db := config.InitDB()
	defer db.Close()

	query := `insert into users(username,email,phone,password) values($1,$2,$3,$4) returning id`

	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}

	// hash password
	hashedPassword, errHash := Hash(payload.Password)
	if errHash != nil {
		log.Fatal(errHash)
	}

	payload.Password = string(hashedPassword)

	var lastId uint64

	queryErr := stmt.QueryRow(payload.Username, payload.Email, payload.Phone, payload.Password).Scan(&lastId)
	if queryErr != nil {
		log.Fatal(queryErr)
	}

	u.Username = payload.Username
	u.Email = payload.Email
	u.Phone = payload.Phone
	u.Password = payload.Password

	u.ID = lastId

	return u, nil
}

// User driver
func SaveDriver(payload schemas.RegisterGuide) (schemas.UserSchema, error) {
	var u schemas.UserSchema

	db := config.InitDB()
	defer db.Close()

	query := `insert into users(username,email,phone,password,image,region,role,language) values($1,$2,$3,$4,$5,$6,$7,$8) returning id`

	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}

	// hash password
	hashedPassword, errHash := Hash(payload.Register.Password)
	if errHash != nil {
		log.Fatal(errHash)
	}

	payload.Register.Password = string(hashedPassword)

	var language string
	var lastId uint64

	if payload.Role == "translator" {
		language = payload.RegisterTranslator.Language
	}

	queryErr := stmt.QueryRow(payload.Register.Username, payload.Register.Email, payload.Register.Phone, payload.Register.Password, payload.Image, payload.Region, payload.Role, language).Scan(&lastId)
	if queryErr != nil {
		log.Fatal(queryErr)
	}

	u.Username = payload.Register.Username
	u.Email = payload.Register.Email
	u.Phone = payload.Register.Phone
	u.Password = payload.Register.Password
	u.Image = payload.Image
	u.Region = payload.Region
	u.Role = payload.Role
	u.ID = lastId

	return u, nil
}

func CheckEmailExists(email string) int {

	var count int

	db := config.InitDB()
	defer db.Close()

	query := `select count(id) from users where email = $1`

	db.QueryRow(query, email).Scan(&count)

	return count
}

func CheckPhoneExists(phone string) int {
	var count int

	db := config.InitDB()
	defer db.Close()

	query := `select count(id) from users where phone = $1`

	db.QueryRow(query, phone).Scan(&count)

	return count
}

func VerifyLogin(email, password string) (schemas.UserSchema, error) {
	var user schemas.UserSchema

	db := config.InitDB()
	defer db.Close()

	query := `select id,username,email,password,phone,avatar,role from users where email = $1`

	row := db.QueryRow(query, email)
	row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Phone, &user.Avatar, &user.Role)

	// check password
	err := CheckPassword(user.Password, password)

	return user, err
}

func FindUserByEmail(email string) (user schemas.UserSchema) {
	var u schemas.UserSchema

	db := config.InitDB()
	defer db.Close()

	query := `select id,username from users where email = $1`
	row := db.QueryRow(query, email)
	row.Scan(&u.ID, &u.Username)

	return u
}

func ConfirmEmailModel(id string) error {
	db := config.InitDB()
	defer db.Close()

	query := `update confirmation_users set activated = $1 where id = $2`

	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}

	stmt.Exec(true, id)

	return nil
}

func GenerateResendExpired() int64 {
	return int64(time.Now().Unix()) + 300 // add 5 minute
}

func SaveResendEmail(id string, user_id uint64) {
	db := config.InitDB()
	defer db.Close()

	query := `update confirmation_users set id = $1, resend_expired = $2 where user_id = $3`

	stmt, _ := db.Prepare(query)

	// add 5 minute resend expired
	resendExpired := GenerateResendExpired()
	stmt.Exec(id, resendExpired, user_id)

	return
}

func GetUsers() []ListUser {
	// define list user
	users := []ListUser{}

	db := config.InitDB()
	defer db.Close()

	query := `select id,username,email,phone,avatar,role,created_at,updated_at,image,region,language from users`

	rows, _ := db.Query(query)
	for rows.Next() {
		var u ListUser

		rows.Scan(&u.Id, &u.Username, &u.Email, &u.Phone, &u.Avatar, &u.Role, &u.Created_At, &u.Updated_At, &u.Image, &u.Region, &u.Language)
		users = append(users, u)
	}

	return users
}
