package helper

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"gopkg.in/gomail.v2"
)

func SendEmail() {
	from := "notification@e-mongolia.mn"
	to := "891834062934@e-mongolia.mn"
	subject := "Test"
	body := "<html><body><h1>Hello!</h1></body></html>" // Replace with the HTML content of your email

	// Compose the email message
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	// Load the CA bundle file
	caCertFile := "./cert/My_CA_Bundle.ca-bundle"
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		panic(err)
	}

	// Create a certificate pool and add the CA certificate
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	// Create a new SMTP sender
	sender := gomail.NewDialer("mail.e-mongolia.mn", 465, from, "d&E:n\\!TD6U=\\whw")
	sender.TLSConfig = &tls.Config{RootCAs: certPool, InsecureSkipVerify: true}

	// Send the email
	if err := sender.DialAndSend(msg); err != nil {
		panic(err)
	}

	fmt.Println("Email sent successfully!")
}
