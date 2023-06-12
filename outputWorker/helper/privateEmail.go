package helper

import(
	"fmt"
)

func SendMail(){
	civilID := "your_civil_id" // Replace with the actual value
	sendMsg := "Your HTML content"   // Replace with your HTML content

	// Get private email from Redis
	// privateEmail, err := redisLocal.HGet(fmt.Sprintf("users:%s", civilID), "email").Result()
	// if err != nil {
	// 	panic(err)
	// }

	// Compose the email message
	msg := []byte(fmt.Sprintf("Subject: Иргэн танд мэдээлэл хүргэж байна\r\n"+
		"To: %s\r\n"+
		"Content-Type: text/html; charset=utf-8\r\n"+
		"Content-Transfer-Encoding: quoted-printable\r\n"+
		"\r\n"+
		"%s\r\n", privateEmail, sendMsg))

	// Create the SMTP client
	auth := smtp.PlainAuth("", AWS_SES_USER, AWS_SES_PASSWORD, AWS_SMTP)
	addr := fmt.Sprintf("%s:%d", AWS_SMTP, 465)
	host, _, _ := net.SplitHostPort(addr)
	from := mail.Address{Name: "", Address: FROM_EMAIL}
	to := mail.Address{Name: "", Address: privateEmail}

	// Establish a TLS connection to the SMTP server
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: host})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Create the SMTP client and send the email
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		panic(err)
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		panic(err)
	}

	if err = client.Mail(from.Address); err != nil {
		panic(err)
	}

	if err = client.Rcpt(to.Address); err != nil {
		panic(err)
	}

	w, err := client.Data()
	if err != nil {
		panic(err)
	}
	defer w.Close()

	_, err = w.Write(msg)
	if err != nil {
		panic(err)
	}

	fmt.Println("Private email sent successfully.")
}