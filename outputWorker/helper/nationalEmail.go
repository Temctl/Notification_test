package helper

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"github.com/Temctl/E-Notification/util/model"
	"gopkg.in/gomail.v2"
)

func AttentionNatEmail(civilId string, content string, notificationType model.NotificationType) int {
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

	return sendNatEmail(civilId, msg)
}

func RegularNatEmail(civilId string, content string) int {

	// TODO regular email html

	return sendNatEmail(civilId, content)
}

func sendNatEmail(civilId string, content string) int {
	from := "notification@e-mongolia.mn"
	// to := "891834062934@e-mongolia.mn"
	to := civilId + "@e-mongolia.mn"
	subject := "Test"
	body := content

	// Compose the email message
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	// Load the CA bundle file
	caCertFile := "./cert/My_CA_Bundle.ca-bundle"
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	// Create a certificate pool and add the CA certificate
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	// Create a new SMTP sender
	sender := gomail.NewDialer("mail.e-mongolia.mn", 465, from, "d&E:n\\!TD6U=\\whw")
	sender.TLSConfig = &tls.Config{RootCAs: certPool, InsecureSkipVerify: true}

	// Send the email
	if err := sender.DialAndSend(msg); err != nil {
		fmt.Println(err)
		return 0
	}

	fmt.Println("Email sent successfully!")
	return 1
}
