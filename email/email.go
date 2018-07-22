package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
)

//Mail type
type Mail struct {
	To      string
	subject string
	body    string
}

type smtpServer struct {
	host string
	port string
}

func (s *smtpServer) serverName() string {
	return s.host + ":" + s.port
}

func (mail *Mail) parseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)

	if err != nil {
		log.Fatal(err)
		return err
	}

	buf := new(bytes.Buffer)

	if err = t.Execute(buf, data); err != nil {
		log.Fatal(err)
		return err
	}
	mail.body = buf.String()

	return nil

}

func (mail *Mail) buildConfirmationMessage(data interface{}) (string, error) {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", "")
	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	err := mail.parseTemplate("confirmation_template.html", data)

	if err != nil {
		log.Println(err)
		return "", err
	}
	message += "\r\n" + mail.body

	return message, nil
}

func (mail *Mail) buildBidNotificationMessage(data interface{}) (string, error) {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", "")
	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	err := mail.parseTemplate("bid-alert_template.html", data)

	if err != nil {
		log.Println(err)
		return "", err
	}

	message += "\r\n" + mail.body

	return message, nil

}

func (mail *Mail) buildPasswordChangeMessage(data interface{}) (string, error) {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", "")
	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	err := mail.parseTemplate("passwordChange_template.html", data)

	if err != nil {
		log.Println(err)
		return "", err
	}
	message += "\r\n" + mail.body

	return message, nil
}

//SendConfirmationMail send email to new users
func (mail *Mail) SendConfirmationMail(name, url string) error {
	mail.subject = os.Getenv("CONFIRMATION-MAIL-SUBJECT")

	data := map[string]string{
		"name": name,
		"url":  url,
	}
	message, err := mail.buildConfirmationMessage(data)

	if err != nil {
		log.Println(err)
		return err
	}

	smtpServer := smtpServer{host: "smtp.gmail.com", port: "465"}

	auth := smtp.PlainAuth("", os.Getenv("EMAIL"), os.Getenv("EMAIL-PASSWORD"), smtpServer.host)

	// Gmail will reject connection if it's not secure
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.host,
	}

	conn, err := tls.Dial("tcp", smtpServer.serverName(), tlsconfig)

	if err != nil {
		log.Println(err)
		return err
	}

	client, err := smtp.NewClient(conn, smtpServer.host)
	if err != nil {
		log.Println(err)
		return err
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		log.Println(err)
		return err
	}

	// step 2: add all from and to
	if err = client.Mail(os.Getenv("EMAIL")); err != nil {
		log.Println(err)
		return err
	}

	if err = client.Rcpt(mail.To); err != nil {
		log.Println(err)
		return err
	}

	// Data
	w, err := client.Data()
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Println(err)
		return err
	}

	err = w.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	client.Quit()
	return nil

}

//EmailPassword sends a users password to them via email
func (mail *Mail) EmailPassword(password, url string) error {
	mail.subject = "New Password"

	data := map[string]string{
		"password": password,
		"url":      url,
	}

	message, err := mail.buildPasswordChangeMessage(data)

	if err != nil {
		log.Println(err)
		return err
	}

	smtpServer := smtpServer{host: "smtp.gmail.com", port: "456"}

	auth := smtp.PlainAuth("", os.Getenv("EMAIL"), os.Getenv("EMAIL-PASSWORD"), smtpServer.host)

	// Gmail will reject connection if it's not secure
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.host,
	}

	conn, err := tls.Dial("tcp", smtpServer.serverName(), tlsconfig)

	if err != nil {
		log.Println(err)
		return err
	}

	client, err := smtp.NewClient(conn, smtpServer.host)
	if err != nil {
		log.Println(err)
		return err
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		log.Println(err)
		return err
	}

	// step 2: add all from and to
	if err = client.Mail(os.Getenv("EMAIL")); err != nil {
		log.Println(err)
		return err
	}

	if err = client.Rcpt(mail.To); err != nil {
		log.Println(err)
		return err
	}

	// Data
	w, err := client.Data()
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Println(err)
		return err
	}

	err = w.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	client.Quit()
	return nil

}

func (mail *Mail) SendBidAlertMail(name, url string) error {
	mail.subject = os.Getenv("BID-ALERT-SUBJECT")

	data := map[string]string{
		"name": name,
		"url":  url,
	}
	message, err := mail.buildBidNotificationMessage(data)

	if err != nil {
		log.Println(err)
		return err
	}

	smtpServer := smtpServer{host: "smtp.gmail.com", port: "465"}

	auth := smtp.PlainAuth("", os.Getenv("EMAIL"), os.Getenv("EMAIL-PASSWORD"), smtpServer.host)

	// Gmail will reject connection if it's not secure
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.host,
	}

	conn, err := tls.Dial("tcp", smtpServer.serverName(), tlsconfig)

	if err != nil {
		log.Println(err)
		return err
	}

	client, err := smtp.NewClient(conn, smtpServer.host)
	if err != nil {
		log.Println(err)
		return err
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		log.Println(err)
		return err
	}

	// step 2: add all from and to
	if err = client.Mail(os.Getenv("EMAIL")); err != nil {
		log.Println(err)
		return err
	}

	if err = client.Rcpt(mail.To); err != nil {
		log.Println(err)
		return err
	}

	// Data
	w, err := client.Data()
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Println(err)
		return err
	}

	err = w.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	client.Quit()
	return nil

}
