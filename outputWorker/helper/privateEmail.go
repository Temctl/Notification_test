package helper

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	gomail "gopkg.in/mail.v2"

	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/model"
)

func AttentionPrivEmail(civilId string, content string, notificationType model.NotificationType) int {
	htmlUrl := "helper/attention_email.html"
	imgElem := "<img style=\"padding-left: 50px;\" src=\"https://content.e-mongolia.mn/png/drivercard.png\" alt=\"citizenCard\" width=\"300px\">"
	if notificationType == model.NotificationType(rune(4)) {
		imgElem = "<img style=\"padding-left: 50px;\" src=\"https://content.e-mongolia.mn/png/passport1.png\" width=\"220px\" alt=\"foreignCard\">"
	} else if notificationType == model.NotificationType(rune(2)) {
		imgElem = "<img style=\"padding-left: 50px;\" src=\"https://content.e-mongolia.mn/png/citizencard.png\" alt=\"citizenCard\" width=\"300px\">"
	}

	htmlContent, err := os.ReadFile(htmlUrl)
	if err != nil {
		fmt.Println("Error reading HTML file:", err)
		return 0
	}

	baseHTML := string(htmlContent)
	msg := strings.ReplaceAll(baseHTML, "{TEXT_REPLACE}", content)
	msg = strings.ReplaceAll(msg, "{IMG_ELEMENT}", imgElem)

	return SendPrivEmail(civilId, msg)
}

// func RegularPrivEmail(civilId string, content string) int {

// 	// TODO regular html

// 	return sendPrivEmail(civilId, content)
// }

func SendPrivEmail(emailAddress string, content string) int {

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "notification@e-mongolia.mn")

	// Set E-Mail receivers
	m.SetHeader("To", "bilguunee10@yahoo.com")

	// Set E-Mail subject
	m.SetHeader("Subject", "Gomail test subject")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", "<p>"+content+"</p>")

	// Settings for SMTP server
	d := gomail.NewDialer(util.AWS_SMTP, 465, util.AWS_SES_USER, util.AWS_SES_PASSWORD)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true, ServerName: util.AWS_SMTP}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("private email sent success")
	return 1
}
