package main

import (
	"fmt"
	"utils"

	"github.com/go-mail/mail"
	"github.com/google/uuid"
)

func SendEmail(token uuid.UUID, email, name string) error {
	m := mail.NewMessage()
	m.SetHeader("From", utils.LoadEnv("email"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Shopify-Integrator Authentication Token")
	m.SetBody("text/html", fmt.Sprintf("<h2>Hi <b>%s</b>,</h2> <br/> <h3>Token: %v</h3>", name, token.String()))
	d := mail.NewDialer("smtp.gmail.com", 587, utils.LoadEnv("email"), utils.LoadEnv("email_psw"))

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
