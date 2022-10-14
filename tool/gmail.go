package tool

import (
	"crypto/tls"
	"fmt"
	"strings"

	lk "github.com/digisan/logkit"
	gm "gopkg.in/mail.v2"
)

func gmail(chOK chan<- error, sender, pwd, header, body string, recipients ...string) {

	// Create New message
	m := gm.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", sender)

	// Set E-Mail recipients
	m.SetHeader("To", recipients...)

	// Set E-Mail subject
	// header: "Gomail test subject"
	m.SetHeader("Subject", header)

	// Set E-Mail body. You can set plain text or html with text/html
	// body: "This is Gomail test body"
	m.SetBody("text/plain", body)

	// Settings for SMTP server
	d := gm.NewDialer("smtp.gmail.com", 587, sender, pwd)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail, and Tell ok to caller
	err := d.DialAndSend(m)
	chOK <- err

	lk.WarnOnErr("%v", err)
}

// Setting at 'https://www.google.com/settings/security/lesssecureapps'

var (
	from = "wismed.cn@gmail.com"
	pwd  = "We1c0me@GOOGLE"
)

func SetGMail(gmail, password string) error {
	if gmail == "" { // use default when input gmail is empty
		return nil
	}
	if !(strings.HasSuffix(gmail, "@gmail.com") && len(gmail) > len("@gmail.com")) {
		return fmt.Errorf("invalid gmail")
	}
	from, pwd = gmail, password
	return nil
}
