package config

import (
	"go-upcycle_connect-backend/internal"
	"go-upcycle_connect-backend/var/email"
	"os"

	"github.com/wneessen/go-mail"
)

func InitEmail() {
	email.Contact = internal.NewEmailClient(
		os.Getenv("EMAIL_FROM"),
		os.Getenv("EMAIL_HOST"),
		os.Getenv("EMAIL_PORT"),
		os.Getenv("EMAIL_USERNAME"),
		os.Getenv("EMAIL_PASSWORD"),
		mail.SMTPAuthPlain,
		mail.NoTLS,
	)
}
