package main

import (
	"fmt"

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
	d := mail.NewDialer("smtp.gmail.com", 587, "keenan@stock2shop.com", "Re_Ghoul")

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
