package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"html/template"
	"net/smtp"
	"ty/car-prices-master/global"
)

func SendEmailNotification(to string, subject string, data interface{}) error {
	from := global.GVA_CONFIG.Email.From
	nickname := global.GVA_CONFIG.Email.Nickname
	secret := global.GVA_CONFIG.Email.Secret
	host := global.GVA_CONFIG.Email.Host
	port := global.GVA_CONFIG.Email.Port
	isSSL := global.GVA_CONFIG.Email.IsSSL
	tmpl := global.GVA_CONFIG.Notification.Tmpl

	var s = bytes.Buffer{}
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return err
	}
	err = t.Execute(&s, data)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", from, secret, host)
	e := email.NewEmail()
	if nickname != "" {
		e.From = fmt.Sprintf("%s <%s>", nickname, from)
	} else {
		e.From = from
	}
	toList := make([]string, 0)
	toList = append(toList, to)
	e.To = toList
	e.Subject = subject
	e.HTML = s.Bytes()
	hostAddr := fmt.Sprintf("%s:%d", host, port)
	if isSSL {
		err = e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: host})
	} else {
		err = e.Send(hostAddr, auth)
	}
	return err
}
