package common

import (
	"crypto/tls"
	"log"
	"net/smtp"
)

type Result struct {
	Success bool
	Message string
}

func PanicIf(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// 由13269827772@163.com发送
// 支持html 单个发送
func SendMail(to string, body string, subject string) error {
	from := "13269827772@163.com"
	password := "a13115681225"
	host := "smtp.163.com"
	port := "465"
	auth := smtp.PlainAuth("", from, password, host)
	nickname := "wScheduler"
	contentType := "Content-Type: text/html; charset=UTF-8"
	msg := []byte("To: " + to + "\r\n" +
		"From: " + nickname + "<" + from + ">\r\n" +
		"Subject: " + subject + "\r\n" +
		contentType + "\r\n\r\n" +
		body)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", host+":"+port, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	if err = c.Auth(auth); err != nil {
		return err
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	if err = c.Rcpt(to); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()

	return nil
}
