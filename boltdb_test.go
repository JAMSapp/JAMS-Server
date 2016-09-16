package main

import (
	"testing"
)

const ID = 123
const USER = "test1"
const PASS = "hunter2"

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

func TestUserLifecycle(t *testing.T) {
	user := &User{
		Id:       ID,
		Username: USER,
		Password: PASS,
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

	u2, err := db.GetUserById(123)
	if err != nil {
		t.Errorf(err.Error())
	}
	if u2.Id != ID {
		t.Errorf("Id of retrieved user does not match stored user")
	}
	if u2.Username != USER {
		t.Errorf("Username of retrieved user does not match stored user")
	}
	if u2.Password != PASS {
		t.Errorf("Password of retrieved user does not match stored user")
	}

	return
}
