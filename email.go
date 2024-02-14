package main

import (
	"fmt"
	"utils"

	"github.com/go-mail/mail"
	"github.com/google/uuid"
)

func Email(token uuid.UUID, email, name string) error {
	m := mail.NewMessage()
	m.SetHeader("From", utils.LoadEnv("email"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Shopify Integrator Registration")
	m.SetBody("text/html", fmt.Sprintf("Hi %s, <br /><p>Please find your registration token below.</p>Token: <br><textarea style='resize:none' rows='2' cols='40'>%v</textarea>", name, token.String()))
	d := mail.NewDialer("smtp.gmail.com", 587, utils.LoadEnv("email"), utils.LoadEnv("email_psw"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
