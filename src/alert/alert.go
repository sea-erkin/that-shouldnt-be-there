package alert

import (
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"strings"

	"github.com/scorredoira/email"
	"github.com/sea-erkin/that-shouldnt-be-there/src/common"
	"github.com/sea-erkin/that-shouldnt-be-there/src/repo"
)

func print(params ...interface{}) {
	fmt.Println(params)
}

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

func PrepNmapScreenshot(ipPorts []repo.IPPortDb, fileName, screenshotTodoDirectory string) {
	// Prep identified ports to take screenshot of new host ports
	print("Prepping screenshots for ports identified as open")

	var portsToScreenshot = make([]string, 0)
	for _, item := range ipPorts {
		if item.State == "open" {
			url := item.IP + ":" + item.Port
			portsToScreenshot = append(portsToScreenshot, url)
		}
	}

	print("Ports to screenshot:", portsToScreenshot)
	print("File name", fileName)
	print("Screenshot directory: ", screenshotTodoDirectory)

	f, err := os.OpenFile(screenshotTodoDirectory+fileName, os.O_APPEND|os.O_WRONLY, 0666)
	for _, i := range portsToScreenshot {
		_, err = f.WriteString(i + "\n")
	}

	f.Close()
	common.CheckErr(err)
}

func CreateEmailBodyFromAlertableHosts(newHosts []repo.HostDb) string {
	body := "New host changes identified \n ================== \n"
	for _, item := range newHosts {
		body += item.Host + "\t New \n"
	}
	print("Email Body", body)
	return body
}

func CreateEmailBodyFromAlertablePorts(addedPorts []repo.IPPortDb) string {
	body := "New port changes identified \n ================== \n"
	for _, port := range addedPorts {
		body += "New State \t" + port.IP + ":" + port.Port + "\t [" + port.Protocol + "]" + "\t" + port.State + "\n"
	}
	print("Email Body", body)
	return body
}
