package email

import (
	"fmt"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/var/email"
	ht "html/template"
	tt "text/template"

	"github.com/wneessen/go-mail"
)

type User struct {
	Firstname string
	Lastname  string
	EmailAddr string
}

func SendEmail(from, to, subject, body string, attachments ...*mail.File) error {
	message := mail.NewMsg()
	if err := message.From(from); err != nil {
		log.Fatal(fmt.Errorf("failed to set From address: %s", err))
		return err
	}
	if err := message.To(to); err != nil {
		log.Fatal(fmt.Errorf("failed to set To address: %s", err))
		return err
	}

	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, body)

	if len(attachments) > 0 {
		message.SetAttachments(attachments)
	}

	if err := email.Contact.Client.DialAndSend(message); err != nil {
		log.Fatal(fmt.Errorf("failed to deliver mail: %s", err))
		return err
	}

	log.Info("Email sent successfully")
	return nil
}

func SendEmails(fromName, fromAddr string, users []User, subject string, textTpl *tt.Template, htmlTpl *ht.Template) error {
	messages, err := prepareBulkMessages(fromName, fromAddr, users, subject, textTpl, htmlTpl)
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return fmt.Errorf("no valid messages to send")
	}

	if err := email.Contact.Client.DialAndSend(messages...); err != nil {
		log.Fatal(fmt.Errorf("failed to deliver bulk mail: %s", err))
		return err
	}

	log.Info(fmt.Sprintf("Bulk mailing successfully delivered to %d recipients", len(messages)))
	return nil
}

func prepareBulkMessages(fromName, fromAddr string, users []User, subject string, textTpl *tt.Template, htmlTpl *ht.Template) ([]*mail.Msg, error) {
	var messages []*mail.Msg
	for _, user := range users {
		message, err := createMessageForUser(fromName, fromAddr, user, subject, textTpl, htmlTpl)
		if err != nil {
			continue
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func createMessageForUser(fromName, fromAddr string, user User, subject string, textTpl *tt.Template, htmlTpl *ht.Template) (*mail.Msg, error) {
	message := mail.NewMsg()

	if err := message.FromFormat(fromName, fromAddr); err != nil {
		log.Fatal(fmt.Errorf("failed to set formatted FROM address: %s", err))
		return nil, err
	}
	if err := message.AddToFormat(fmt.Sprintf("%s %s", user.Firstname, user.Lastname), user.EmailAddr); err != nil {
		log.Fatal(fmt.Errorf("failed to set formatted TO address: %s", err))
		return nil, err
	}

	message.SetMessageID()
	message.SetDate()
	message.SetBulk()
	message.Subject(subject)

	if textTpl != nil {
		if err := message.SetBodyTextTemplate(textTpl, user); err != nil {
			log.Fatal(fmt.Errorf("failed to add text template to mail body: %s", err))
			return nil, err
		}
	}

	if htmlTpl != nil {
		if err := message.AddAlternativeHTMLTemplate(htmlTpl, user); err != nil {
			log.Fatal(fmt.Errorf("failed to add HTML template to mail body: %s", err))
			return nil, err
		}
	}

	return message, nil
}
