package main

import (
	"fmt"
	"utils"

	"github.com/go-mail/mail"
)

func SendEmail(token, email, name string) error {
	m := mail.NewMessage()
	m.SetHeader("From", "keenan@stock2shop.com")
	m.SetHeader("To", email)
	// m.SetAddressHeader("Cc", "oliver.doe@example.com", "Oliver")
	m.SetHeader("Subject", "Shopify-Integrator Authentication Token")
	m.SetBody("text/html", fmt.Sprintf("Hi <b>%s</b>, <br> Token: %s", name, token))
	// m.Attach("lolcat.jpg")
	d := mail.NewDialer("smtp.gmail.com", 587, utils.LoadEnv("email"), utils.LoadEnv("email_psw"))

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
