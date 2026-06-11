package internal

import (
	"fmt"
	"go-upcycle_connect-backend/utils/log"
	"strconv"

	"github.com/wneessen/go-mail"
)

type Email struct {
	Client *mail.Client
	From   string
}

func NewEmailClient(from, host, port, username, password string, authType mail.SMTPAuthType, tlsPolicy mail.TLSPolicy) Email {
	var intPort, errorToConvert = strconv.Atoi(port)

	if errorToConvert != nil {
		log.Fatal(errorToConvert)
	}

	client, err := mail.NewClient(host,
		mail.WithSMTPAuth(authType),
		mail.WithTLSPortPolicy(tlsPolicy),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithPort(intPort),
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Info(fmt.Sprintf("(CONFIG) Email Client Initialized %s", from))

	return Email{Client: client, From: from}
}
