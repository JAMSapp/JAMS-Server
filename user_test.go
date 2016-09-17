package main

import (
	"testing"
)

func TestLifecycle(t *testing.T) {
	user := &User{Id: 1234, Username: "user", Password: "hunter2"}
	if err := user.Save(); err != nil {
		t.Errorf(err.Error())
	}

	if err := user.Delete(); err != nil {
		t.Errorf(err.Error())
	}

	user, err := db.GetUserById(1234)
	if err != ErrUserNotFound {
		t.Errorf("Database found deleted user")
	}
}
