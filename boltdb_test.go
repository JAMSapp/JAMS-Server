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
	uuid.Init() // Must init before V1 uuids.
	mes := NewMessage("test message body")
	err = db.SaveMessage(mes)
	if err != nil {
		t.Errorf(err.Error())
	}

	user := &User{Id: ID, Username: USER, Password: PASS}

	err = db.AddUnreadMessage(user, mes)
	if err != nil {
		t.Errorf(err.Error())
	}

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

	user, err = db.GetUserById(ID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if user == nil {
		t.Errorf("User should not be nil.")
		return
	}
	if user.Id == ID {
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
	//uuid.Init() // Must init before V1 uuids.
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
