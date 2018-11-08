// SMTP is Simple Mail Transfer Protocol, there are completing sending mail's client and local mail server with golang.

/*
 Common email addresses and ports:
 1. qq.com:465/587
 2. 163.com:465/587
 3. 126.com:465/587
*/

package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

// Send SMTP mail, it's mail protocol always.
func SendSMTPMail(user, password, host, to, subject, body, mailType string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var contentType string
	if mailType == "html" {
		contentType = "Content-Type: text/" + mailType + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + "\r\n" + contentType + "\r\n\r\n" + body)
	sendTo := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, sendTo, msg)
	return err
}

func main() {
	user := "yang**@yun*.com"
	password := "***"
	host := "smtp.exmail.qq.com:25"
	to := "397685131@qq.com"
	subject := "使用Golang发送邮件"
	body := `
		<html>
		<body>
		<h3>
		"Test send to email"
		</h3>
		</body>
		</html>
		`
	fmt.Println("send email")
	err := SendSMTPMail(user, password, host, to, subject, body, "html")
	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("Send mail success!")
	}

}
