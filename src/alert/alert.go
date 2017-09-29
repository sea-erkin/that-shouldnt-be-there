package alert

import (
	"log"
	"net/smtp"
	"strings"
)

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
