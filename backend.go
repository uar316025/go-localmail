package localmail

type Backend struct {
	Users map[string]*User
}

func NewBackend() *Backend {
	return &Backend{
		Users: map[string]*User{},
	}
}

func (back *Backend) getOrCreate(username string) *User {
	user, ok := back.Users[username]
	if !ok {
		user = NewUser(username)
		back.Users[username] = user
	}
	return user
}
