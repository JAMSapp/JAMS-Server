package main

import (
	"testing"
)

func TestMarshal(t *testing.T) {
	user := &User{
		Id:       123,
		Username: "farts",
		Password: "farts",
	}
	t.Logf("%s", MarshalUser(user))
	return
}

func TestSaveUser(t *testing.T) {
	user := &User{
		Id:       123,
		Username: "farts",
		Password: "farts",
	}

	db, err := BoltDBOpen()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	err = db.SaveUser(user)
	if err != nil {
		t.Errorf(err.Error())
	}
	return
}
