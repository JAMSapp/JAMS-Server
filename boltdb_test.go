package main

import (
	"testing"
)

func TestBoltDBOpen(t *testing.T) {
	db, err := BoltDBOpen("my.db")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	// Try and open another.
	_, err = BoltDBOpen("my.db")
	if err == nil {
		t.Errorf("boltdb: error expected with opening dir")
	}
}

func TestSaveUser(t *testing.T) {
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

	err = db.SaveUser(user)
	if err != nil {
		t.Errorf(err.Error())
	}
	return
}
