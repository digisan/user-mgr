package tool

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/digisan/gotk/strs"
	lk "github.com/digisan/logkit"
	gm "gopkg.in/mail.v2"
)

func gmail(chOK chan<- error, from, pwd, header, body string, to ...string) {

	// Create New message
	m := gm.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", from)

	// Set E-Mail receivers
	receivers := to
	m.SetHeader("To", receivers...)

	// Set E-Mail subject
	// header: "Gomail test subject"
	m.SetHeader("Subject", header)

	// Set E-Mail body. You can set plain text or html with text/html
	// body: "This is Gomail test body"
	m.SetBody("text/plain", body)

	// Settings for SMTP server
	d := gm.NewDialer("smtp.gmail.com", 587, from, pwd)

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

func genCode(email string) string {
	key := []byte(fmt.Sprintf("%d", time.Now().UnixNano())[3:19])
	return strs.Maxlen(fmt.Sprintf("%06x", Encrypt(email, key)), 6)
}

func SendCode(ctx context.Context, to string, timeout time.Duration) (string, error) {

	header := "Sign-Up Verification Code"

	code := genCode(to)
	// fmt.Println(code)
	body := fmt.Sprintf("verification code: %s\n", code)

	chOK := make(chan error)
	go gmail(chOK, from, pwd, header, body, to)

	select {
	case <-ctx.Done():
		// fmt.Printf("%v\n", "out cancelled")
		return "", fmt.Errorf("out cancelled")

	case <-time.After(timeout):
		// fmt.Printf("%v\n", "time out")
		return "", fmt.Errorf("time out")

	case err := <-chOK:
		// fmt.Printf("%v\n", err)
		return code, err
	}
}
