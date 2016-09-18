package main

import (
	"testing"
)

const ID = 123
const USER = "test1"
const PASS = "hunter2"

func TestBoltDBOpen(t *testing.T) {
	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	// Try and open another.
	_, err = BoltDBOpen(DBFILE)
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

	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	err = db.SaveUser(user)
	if err != nil {
		t.Errorf(err.Error())
	}

	user, err = db.GetUserById(ID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if user == nil {
		t.Errorf("User should not be nil.")
		return
	}
	if user.Id != ID {
		t.Errorf("Id of retrieved user does not match stored user")
	}
	if user.Username != USER {
		t.Errorf("Username of retrieved user does not match stored user")
	}
	if user.Password != PASS {
		t.Errorf("Password of retrieved user does not match stored user")
	}

	// TODO: Check the users result.
	_, err = db.GetUsers()
	if err != nil {
		t.Errorf(err.Error())
	}

	err = db.DeleteUser(user)
	if err != nil {
		t.Errorf(err.Error())
	}

	user, err = db.GetUserById(ID)
	if err != ErrUserNotFound {
		t.Errorf("Found deleted user")
	}
	if user != nil {
		t.Errorf("GetUserById should have returned nill")
	}

	return
}

func TestMarshalUser(t *testing.T) {
	user := &User{
		Id:       123,
		Username: "farts",
		Password: "farts",
	}
	t.Logf("%s", MarshalUser(user))
	return
}
