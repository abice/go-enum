package prompt_test

import (
	"net/mail"

	"github.com/Bowery/prompt"
	"golang.org/x/crypto/bcrypt"
)

func ExamplePassword() ([]byte, error) {
	clear, err := prompt.Password("Password")
	if err != nil {
		return nil, err
	}

	return bcrypt.GenerateFromPassword([]byte(clear), bcrypt.DefaultCost)
}

func ExampleBasic() (string, error) {
	email, err := prompt.Basic("Email", true)
	if err != nil {
		return "", err
	}

	_, err = mail.ParseAddress(email)
	return email, err
}
