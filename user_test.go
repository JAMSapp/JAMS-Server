package main

import (
	"testing"
)

func TestMarshalUser(t *testing.T) {
	user := &User{
		Id:       123,
		Username: "farts",
		Password: "farts",
	}
	t.Logf("%s", MarshalUser(user))
	return
}

func TestSave(t *testing.T) {
	user := &User{
		Id:       123,
		Username: "farts",
		Password: "farts",
	}
	db, err := BoltDBOpen("my.db")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	err = user.Save(db)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	err = user.Save(nil)
	if err == nil {
		t.Errorf("expected error on saving against nil dbconn")
	}
	return
}
