package localmail

import (
	"errors"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
	"log"
)

type imapBackend struct {
	*Backend
}

func (back *Backend) IMAP() *imapBackend {
	return &imapBackend{back}
}

func (back *imapBackend) Login(connInfo *imap.ConnInfo, username, password string) (backend.User, error) {
	if 5 >= len(password) {
		return nil, errors.New("incorrect password, should be at least 6 chars")
	}

	user, ok := back.Users[username]
	if !ok {
		log.Printf("[IMAP] mailbox not found for '%s'\n", username)

		return nil, errors.New("user not found")
	}
	log.Printf("[IMAP] get mailbox for '%s'\n", username)
	return user, nil
}
