package email

import "gopkg.in/gomail.v2"

type MailBox struct {
	dialer *gomail.Dialer
}

func NewMailBox(smtpServerAddr string, smtpPort int, login string, password string) MailBox {
	d := gomail.NewDialer(smtpServerAddr, smtpPort, login, password)
	mb := MailBox{
		dialer: d,
	}
	return mb
}

func (mb MailBox) SendEmail(fromAddr string, toAddr string, emailSubject string, emailText string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", fromAddr)
	msg.SetHeader("To", toAddr)
	msg.SetHeader("Subject", emailSubject)
	msg.SetBody("text/plain", emailText)
	err := mb.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
