package alert

import (
	"log"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/scorredoira/email"
)

func SendMailAttachment(body, fromEmail, password, filePath, emailHost, emailPort string, to []string) {
	// compose the message
	m := email.NewMessage("[TSBT]", body)
	m.From = mail.Address{Name: "From", Address: fromEmail}
	m.To = to

	// add attachments
	if err := m.Attach(filePath); err != nil {
		log.Fatal(err)
	}

	// send it
	auth := smtp.PlainAuth("", fromEmail, password, emailHost)
	if err := email.Send(emailHost+":"+emailPort, auth, m); err != nil {
		log.Fatal(err)
	}
}

func SendMail(body, fromEmail, password, emailHost, emailPort string, to []string) {
	from := fromEmail
	pass := password

	msg := "From: " + from + "\n" +
		"To: " + strings.Join(to, ";") + "\n" +
		"Subject: [TSBT]\n\n" +
		body

	err := smtp.SendMail(emailHost+":"+emailPort,
		smtp.PlainAuth("TSBT", from, pass, emailHost),
		from, to, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}
