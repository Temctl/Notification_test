package helper

import (
	"fmt"
	"net/mail"
	"net/smtp"
)

func SendEmail() {
	from := mail.Address{Name: "notification@e-mongolia.mn", Address: "notification@e-mongolia.mn"}
	to := mail.Address{Name: "", Address: "891834062934@e-mongolia.mn"} // Replace with the recipient's email address
	subject := "Иргэн танд мэдээлэл хүргэж байна"
	body := "<html><body><h1>Hello!</h1></body></html>" // Replace with the HTML content of your email

	// Compose the email
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject
	headers["Content-Type"] = "text/html; charset=utf-8"

	var msg string
	for key, value := range headers {
		msg += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	msg += "\r\n" + body

	// Set up the SMTP client
	auth := smtp.PlainAuth("", from.Address, "d&E:n\\!TD6U=\\whw", "mail.e-mongolia.mn")
	smtpServer := "mail.e-mongolia.mn:465" // Replace with the SMTP server and port

	// Send the email
	err := smtp.SendMail(smtpServer, auth, from.Address, []string{to.Address}, []byte(msg))
	if err != nil {
		fmt.Println("Error sending email:", err)
		return
	}

	fmt.Println("Email sent successfully!")
}
