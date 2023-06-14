package helper

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strings"

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

	return sendPrivEmail(civilId, msg)
}

func sendPrivEmail(emailAddress string, content string) int {

	// privateEmail := "temuulen@ema.gov.mn"
	privateEmail := "utemuka@gmail.com"

	// Compose the email message
	msg := []byte(fmt.Sprintf("Subject: Иргэн танд мэдээлэл хүргэж байна\r\n"+
		"To: %s\r\n"+
		"Content-Type: text/html; charset=utf-8\r\n"+
		"Content-Transfer-Encoding: quoted-printable\r\n"+
		"\r\n"+
		"%s\r\n", privateEmail, content))

	// Create the SMTP client
	auth := smtp.PlainAuth("", util.AWS_SES_USER, util.AWS_SES_PASSWORD, util.AWS_SMTP)
	addr := fmt.Sprintf("%s:%d", util.AWS_SMTP, 465)
	host, _, _ := net.SplitHostPort(addr)
	// from := mail.Address{Name: "", Address: util.FROM_EMAIL}
	to := mail.Address{Name: "", Address: privateEmail}

	// Establish a TLS connection to the SMTP server
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: host})
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer conn.Close()

	// Create the SMTP client and send the email
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		fmt.Println(err)
		return 0
	}

	if err = client.Mail("notification@e-mongolia.mn"); err != nil {
		fmt.Println(err)
		return 0
	}

	if err = client.Rcpt(to.Address); err != nil {
		fmt.Println(err)
		return 0
	}

	w, err := client.Data()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer w.Close()

	_, err = w.Write(msg)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	fmt.Println("Private email sent successfully.")
	return 1
}
