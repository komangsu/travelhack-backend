package libs

import (
	"bytes"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strconv"
	"text/template"
)

func ParseTemplate(templateName string, data interface{}) (string, error) {
	// html to byte
	t, err := template.ParseFiles(templateName)
	if err != nil {
		return "", err
	}

	// render struct data to html file
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func SendEmail(to string, subject string, data interface{}, templateFile string) error {
	var (
		SMTP_HOST     = os.Getenv("SMTP_HOST")
		SMTP_PORT     = os.Getenv("SMTP_PORT")
		AUTH_USER     = os.Getenv("AUTH_USER")
		AUTH_PASSWORD = os.Getenv("AUTH_PASSWORD")
	)

	// change port into int
	PORT, _ := strconv.Atoi(SMTP_PORT)

	result, _ := ParseTemplate(templateFile, data)
	m := gomail.NewMessage()
	m.SetHeader("From", "travel@travelhack.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", result)

	dialer := gomail.NewDialer(
		SMTP_HOST,
		PORT,
		AUTH_USER,
		AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(m)
	if err != nil {
		panic(err)
	}
	return err
}

func SendEmailVerification(to string, data interface{}) {
	var err error

	template := "./templates/email/email-confirm.html"
	subject := "Email Verification"

	err = SendEmail(to, subject, data, template)
	if err == nil {
		log.Println("Success send email.")
	} else {
		log.Println(err)
	}
}
