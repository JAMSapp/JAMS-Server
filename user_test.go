package main

import (
	"testing"
)

func TestUserLifecycle(t *testing.T) {
	var err error
	db, err = BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf("On db open: %s", err.Error())
		return
	}
	defer db.Close()

	user, err := NewUser("username", "password")
	if err != nil {
		t.Errorf("On new user: %s", err.Error())
	}

	if err = user.Save(); err != nil {
		t.Errorf("On user save: %s", err.Error())
	}

	if err = user.Delete(); err != nil {
		t.Errorf("On user delete: %s", err.Error())
	}

	_, err = db.GetUserByUsername("username")
	if err != ErrUserNotFound {
		t.Errorf("Database found deleted user")
	}
}

func TestMessageSend(t *testing.T) {
	var err error
	db, err = BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf("On db open: %s", err.Error())
		return
	}
	defer db.Close()

	u1, _ := NewUser("user1", "password")
	m1 := NewMessage(MBODY)
	if err = u1.SendMessage(m1); err != nil {
		t.Errorf("On message send: %s", err.Error())
	}
}
