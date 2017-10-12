package alert

import (
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/scorredoira/email"
	"github.com/sea-erkin/that-shouldnt-be-there/src/common"
	"github.com/sea-erkin/that-shouldnt-be-there/src/repo"
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

func PrepNmapScreenshot(ipPorts []repo.IPPortDb, fileName, screenshotTodoDirectory string) {
	// Prep identified ports to take screenshot of new host ports

	print("Prepping screenshots")

	var portsToScreenshot = make([]string, 0)
	for _, item := range ipPorts {
		url := item.IP + ":" + item.Port
		portsToScreenshot = append(portsToScreenshot, url)
	}

	print("Ports to screenshot:", portsToScreenshot)
	// Create file in nmap todo directory with the same params
	newLineString := strings.Join(portsToScreenshot, "\n")
	d1 := []byte(newLineString)

	print("File name", fileName)
	print("Screenshot directory: ", screenshotTodoDirectory)
	err := ioutil.WriteFile(screenshotTodoDirectory+fileName, d1, 0644)
	common.CheckErr(err)
}

func CreateEmailBodyFromAlertableHosts(newHosts []repo.HostDb, missingHosts []repo.HostDb) string {
	body := "New host changes identified \n ================== \n"
	for _, item := range newHosts {
		body += item.Host + "\t New \n"
	}
	for _, item := range missingHosts {
		body += item.Host + "\t Missing \n"
	}
	print("Email Body", body)
	return body
}

func CreateEmailBodyFromAlertablePorts(addedPorts []repo.IPPortDb, missingPorts []repo.IPPortDb) string {
	body := "New port changes identified \n ================== \n"
	for _, port := range addedPorts {
		body += "New State \t" + port.IP + ":" + port.Port + "\t [" + port.Protocol + "]" + "\t" + port.State + "\n"
	}
	for _, port := range missingPorts {
		body += "Previous State \t" + port.IP + ":" + port.Port + "\t [" + port.Protocol + "]" + "\t" + port.State + "\n"
	}
	print("Email Body", body)
	return body
}
