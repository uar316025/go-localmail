package localmail

import (
	"bytes"
	"errors"
	"github.com/emersion/go-smtp"
	"io"
	"log"
	"strings"
	"time"
)

type smtpSession struct {
	back *Backend
	user *User
	date time.Time
}

func newSmtpSession(back *Backend) *smtpSession {
	sess := &smtpSession{
		back: back,
		date: time.Now(),
	}
	sess.Reset()
	return sess
}

func (sess *smtpSession) Reset() {
}

func (sess *smtpSession) Logout() error {
	return nil
}

func (sess *smtpSession) Mail(from string, opts smtp.MailOptions) error {
	return nil
}

func (sess *smtpSession) Rcpt(to string) error {
	if idx := strings.Index(to, "@"); idx != 0 {
		to = to[:idx]
	}
	sess.user = sess.back.getOrCreate(to)
	return nil
}

func (sess *smtpSession) Data(reader io.Reader) error {
	if sess.user == nil {
		return errors.New("no recipient defined")
	}

	log.Printf("[SMTP] recieced mail for '%s'\n", sess.user.username)

	mbox, err := sess.user.GetMailbox("INBOX")
	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)
	_, _ = buffer.ReadFrom(reader)

	return mbox.CreateMessage(
		[]string{"\\Recent", "\\Unseen"},
		time.Now(),
		buffer,
	)

}
