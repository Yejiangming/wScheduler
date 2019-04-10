package common

import (
	"log"
	"net/smtp"
	"strings"
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
// 支持html 支持向多个邮箱发送
func SendMail(to []string, msg string, subject string) error {
	username := "13269827772@163.com"
	password := "a13115681225"
	host := "smtp.163.com"
	auth := smtp.PlainAuth("", username, password, host)
	nickname := "wScheduler"
	contentType := "Content-Type: text/html; charset=UTF-8"
	body := []byte("To: " + strings.Join(to, ",") + "\r\n" +
		"From: " + nickname + "<" + username + ">\r\n" +
		"Subject: " + subject + "\r\n" +
		contentType + "\r\n\r\n" +
		msg)
	err := smtp.SendMail("smtp.163.com:25", auth, username, to, body)
	return err
}
