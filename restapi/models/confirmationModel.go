package models

import (
	"restapi/config"
	"time"
)

type Confirmation struct {
	Id             string `json:"id"`
	Activated      bool   `json:"activated"`
	Resend_expired int64  `json:"resend_expired"`
	User_id        uint64 `json:"user_id"`
}

func SaveConfirmation(id string, user_id uint64) error {
	var confirmation Confirmation

	db := config.InitDB()
	defer db.Close()

	query := `insert into confirmation_users(id,user_id,resend_expired) values($1,$2,$3)`

	stmt, _ := db.Prepare(query)

	confirmation.Resend_expired = time.Now().Add(time.Minute * 15).Unix()
	stmt.Query(id, user_id, confirmation.Resend_expired)
	return nil
}

func FindConfirmation(user_id uint64) (Confirmation, error) {
	var confirmation Confirmation

	db := config.InitDB()
	defer db.Close()

	query := `select * from confirmation_users where user_id = $1`

	row := db.QueryRow(query, user_id)
	row.Scan(&confirmation.Id, &confirmation.Activated, &confirmation.Resend_expired, &confirmation.User_id)

	return confirmation, nil
}

func CheckActivatedById(id string) Confirmation {
	var confirmation Confirmation

	db := config.InitDB()
	defer db.Close()

	query := `select activated,user_id from confirmation_users where id = $1`

	row := db.QueryRow(query, id)
	row.Scan(&confirmation.Activated, &confirmation.User_id)

	return confirmation
}
