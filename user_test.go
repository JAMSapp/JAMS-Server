package main

import (
	"testing"
)

func TestSave(t *testing.T) {
	user := &User{Id: 1234, Username: "user", Password: "hunter2"}
	if err := user.Save(); err != nil {
		t.Errorf(err.Error())
	}
}
