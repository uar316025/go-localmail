package localmail

import (
	"errors"
	"github.com/emersion/go-smtp"
)

type smtpBackend struct {
	*Backend
}

func (back *smtpBackend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	if 5 >= len(password) {
		return nil, errors.New("incorrect password, should be at least 6 chars")
	}

	return newSmtpSession(back.Backend), nil
}

func (back *smtpBackend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return nil, errors.New("not allowed")
}

func (back *Backend) SMTP() *smtpBackend {
	return &smtpBackend{back}
}
