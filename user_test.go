package main

import (
	"testing"
)

func TestUserLifecycle(t *testing.T) {
	user := &User{Id: "asdf", Username: "user", Password: "hunter2"}
	if err := user.Save(); err != nil {
		t.Errorf(err.Error())
	}

	if err := user.Delete(); err != nil {
		t.Errorf(err.Error())
	}

	user, err := db.GetUserById("asdf")
	if err != ErrUserNotFound {
		t.Errorf("Database found deleted user")
	}
}
