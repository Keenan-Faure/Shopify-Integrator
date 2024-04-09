package main

import (
	"errors"
	"fmt"
	"utils"

	"github.com/go-mail/mail"
	"github.com/google/uuid"
)

func Email(token uuid.UUID, sendEmail bool, email, name string) error {
	valid, err := emailValidation(token, email, name)
	if err != nil {
		return err
	}
	if valid {
		m := mail.NewMessage()
		m.SetHeader("From", utils.LoadEnv("email"))
		m.SetHeader("To", email)
		m.SetHeader("Subject", "Shopify Integrator Registration")
		m.SetBody("text/html", fmt.Sprintf("Hi %s, <br /><p>Please find your registration token below.</p>Token: <br><textarea style='resize:none' rows='2' cols='40'>%v</textarea>", name, token.String()))
		d := mail.NewDialer("smtp.gmail.com", 587, utils.LoadEnv("email"), utils.LoadEnv("email_psw"))
		if sendEmail {
			if err := d.DialAndSend(m); err != nil {
				return err
			}
		}
		return nil
	}
	// should never get here, but you never knw
	return errors.New("invalid email cannot be sent")
}

func emailValidation(token uuid.UUID, email, name string) (bool, error) {
	if email == "" || len(email) == 0 {
		return false, errors.New("invalid email not allowed")
	}
	if name == "" || len(name) == 0 {
		return false, errors.New("invalid name not allowed")
	}
	if token == uuid.Nil {
		return false, errors.New("invalid token not allowed")
	}
	return true, nil
}
