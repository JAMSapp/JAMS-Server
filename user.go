package main

// User represents and user registered with the system.
type User struct {
	Id       int
	Username string
	Password string
}

func (u *User) Save() error {
	return db.SaveUser(u)
}
