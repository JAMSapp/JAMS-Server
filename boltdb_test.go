package main

import (
	"testing"

	"github.com/twinj/uuid"
)

var ID = uuid.NewV1().String()
var USER = "test1"
var PASS = "hunter2"

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

func TestBoltMessageLifecycle(t *testing.T) {
	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	// TODO: Create some sort of init function to handle the calling of any
	// necessary init requirements throughout server.
	mes := NewMessage("test message body")
	err = db.SaveMessage(mes)
	if err != nil {
		t.Errorf(err.Error())
	}

	user := &User{Id: ID, Username: USER, Password: PASS}

	err = db.AddUnreadMessage(user, mes)
	if err != nil {
		t.Errorf("AddUnreadMessage: %s", err.Error())
	}

	_, err = db.GetUnreadMessages(user)
	if err != nil {
		t.Errorf("GetUnreadMessages: %s", err.Error())
	}
	/*
		if len(messages) != 1 {
			t.Errorf("db.GetUnreadMessages returned too many messages")
		}
	*/

	err = db.DeleteMessage(mes)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestBoltUserLifecycle(t *testing.T) {
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

	// Make sure GetByUserId returns correct user
	user2, err := db.GetUserById(user.Id)
	if err != nil {
		t.Errorf(err.Error())
	}
	testUsersEqual(user, user2, t)

	// Make sure GetByUserByUsername returns correct user
	user3, err := db.GetUserByUsername(user.Username)
	if err != nil {
		t.Errorf(err.Error())
	}
	testUsersEqual(user, user3, t)

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
		t.Errorf("GetUserById should have returned nil")
	}

	return
}

func TestBoltMessageSend(t *testing.T) {
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

	// TODO: Create some sort of init function to handle the calling of any
	// necessary init requirements throughout server.
	mes := NewMessage("test message body")
	err = db.SaveMessage(mes)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = db.AddUnreadMessage(user, mes)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestBoltMarshalUser(t *testing.T) {
	user := &User{
		Id:       ID,
		Username: USER,
		Password: PASS,
	}
	t.Logf("%s", MarshalUser(user))
	return
}

func testUsersEqual(u1, u2 *User, t *testing.T) {
	if u2 == nil {
		t.Errorf("User should not be nil.")
		return
	}
	if u1.Id != u2.Id {
		t.Errorf("Id of retrieved user does not match stored user: %s vs %s", u1.Id, u2.Id)
	}
	if u1.Username != u2.Username {
		t.Errorf("Username of retrieved user does not match stored user: %s vs %s", u1.Id, u2.Id)
	}
	if u1.Password != u2.Password {
		t.Errorf("Password of retrieved user does not match stored user: %s vs %s", u1.Id, u2.Id)
	}
}

func testGetUserByUsername(username string, t *testing.T) {
	user, err := db.GetUserByUsername(username)
	if err != nil {
		t.Errorf(err.Error())
	}
	if user == nil {
		t.Errorf("User should not be nil.")
		return
	}
	if user.Id != ID {
		t.Errorf("Id of retrieved user does not match stored user: %s vs %s", user.Id, ID)
	}
	if user.Username != USER {
		t.Errorf("Username of retrieved user does not match stored user")
	}
	if user.Password != PASS {
		t.Errorf("Password of retrieved user does not match stored user")
	}
}
