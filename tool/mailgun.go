package tool

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	cfg "github.com/digisan/go-config"
	lk "github.com/digisan/logkit"
	"github.com/mailgun/mailgun-go/v4"
)

var (
	domain = ""
	apiKey = ""
	sender = ""
	mg     *mailgun.MailgunImpl
)

func init() {

	lk.Log("starting...email")

	if err := cfg.Init("email", false, "mailgun-config.json"); err == nil {
		domain = cfg.Val[string]("domain")
		apiKey = cfg.Val[string]("apikey")
		sender = cfg.Val[string]("sender")
	}

	lk.FailOnErrWhen(len(domain) == 0, "%v", errors.New("email domain is empty"))
	lk.FailOnErrWhen(len(apiKey) == 0, "%v", errors.New("email apiKey is empty"))
	lk.FailOnErrWhen(len(sender) == 0, "%v", errors.New("email sender is empty"))

	mg = mailgun.NewMailgun(domain, apiKey)
}

func SetMailDomain(d string) {
	domain = d
}

func SetMailAPIKey(k string) {
	apiKey = k
}

func SetMailSender(s string) {
	sender = s
}

type sdResult struct {
	recipient string
	msg       string
	id        string
	err       error
}

func send(subject, body string, recipients ...string) chan sdResult {
	var (
		chRst = make(chan sdResult)
		nOK   = int32(0)
	)

	for _, recipient := range recipients {
		go func(recipient string) {

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			// The message object allows you to add attachments and Bcc recipients
			message := mg.NewMessage(sender, subject, body, recipient)

			// Send the message with a 10 second timeout
			if msg, id, err := mg.Send(ctx, message); err != nil {

				lk.Warn("ID: %s Resp: %s Err: %v\n", id, msg, err)
				chRst <- sdResult{
					recipient: recipient,
					msg:       "",
					id:        "",
					err:       err,
				}

			} else {

				lk.Log("ID: %s Resp: %s Err: %v\n", id, msg, err)
				chRst <- sdResult{
					recipient: recipient,
					msg:       msg,
					id:        id,
					err:       nil,
				}
				atomic.AddInt32(&nOK, 1)
			}

		}(recipient)
	}
	return chRst
}

func SendMail(subject, body string, recipients ...string) (OK bool, sent []string, failed []string, errs []error) {
	var (
		timeout = 15 * time.Second
		chRst   = send(subject, body, recipients...)
		nOK     = 0
	)

	select {
	case <-time.After(1 * time.Millisecond):
		for rst := range chRst {
			if rst.err == nil {
				sent = append(sent, rst.recipient)
				nOK++
			} else {
				failed = append(failed, rst.recipient)
				errs = append(errs, rst.err)
			}
			if nOK == len(recipients) {
				close(chRst)
			}
		}
		return nOK == len(recipients), sent, failed, errs

	case <-time.After(timeout):
		errs = append(errs, fmt.Errorf("timeout @%vs", timeout/time.Second))
		return false, nil, nil, errs
	}
}
