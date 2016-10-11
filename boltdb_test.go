package main

import (
	"testing"

	"github.com/twinj/uuid"
)

var ID = uuid.NewV1().String()
var ID2 = uuid.NewV1().String()

var USER = "test1"
var PASS = "hunter2"

var user = &User{Id: ID, Username: USER, Password: PASS}
var userDupeName = &User{Id: ID2, Username: USER, Password: PASS}

var BODY = "This is a test message."

var message = &Message{Id: ID, Body: BODY}

var thread = &Thread{Id: ID, UserIds: []string{ID, ID2}}

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

func TestBoltSaveUser(t *testing.T) {
	// Setup DB
	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	var nilUser *User
	nilUser = nil
	// Test saving a nil user.
	err = db.SaveUser(nilUser)
	if err != ErrUserNil {
		t.Errorf("Saving user with nil object did not return ErrUserNil")
	}

	// Test regular ol' save.
	err = db.SaveUser(user)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Test a user with a duplicate username.
	err = db.SaveUser(userDupeName)
	if err != ErrUsernameAlreadyExists {
		t.Errorf("Saving user with duplicate username did not return ErrUsernameAlreadyExists")
	}
}

func TestBoltGetUserById(t *testing.T) {
	// Setup DB
	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	// Test getting user with blank Id
	u, err := db.GetUserById("")
	if err != ErrUserIdBlank {
		t.Errorf("Getting user with blank id did not return ErrUserIdBlank")
	}
	if u != nil {
		t.Errorf("Getting user with blank username did not return nil user")
	}

	// Test getting user with unknown Id
	u, err = db.GetUserById("1234-1234-1234-1234")
	if err != ErrUserNotFound {
		t.Errorf("Getting user with no known Id did not return ErrUserNotFound")
	}
	if u != nil {
		t.Errorf("Getting user with blank username did not return nil user")
	}

	// Test getting user by a real Id.
	u, err = db.GetUserById(user.Id)
	if err != nil {
		t.Errorf(err.Error())
	}
	if u == nil {
		t.Errorf("GetUserById returned nil user")
	} else {
		testUsersEqual(u, user, t)
	}
}

func TestBoltGetUserByUsername(t *testing.T) {
	// Setup DB
	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	// Test getting user with blank username
	u, err := db.GetUserByUsername("")
	if err != ErrUsernameBlank {
		t.Errorf("Getting user with blank username did not return ErrUsernameBlank")
	}
	if u != nil {
		t.Errorf("Getting user with blank username did not return nil user")
	}

	// Test getting username by unknown username
	u, err = db.GetUserByUsername("1234-1234-1234-1234")
	if err != ErrUserNotFound {
		t.Errorf("Getting user with no known Username did not return ErrUserNotFound")
	}
	if u != nil {
		t.Errorf("Getting user with unknown username did not return nil user")
	}

	// Test getting user by a real username
	u, err = db.GetUserByUsername(user.Username)
	if err != nil {
		t.Errorf(err.Error())
	} else {
		if u == nil {
			t.Errorf("GetUserById returned nil user")
		} else {
			testUsersEqual(u, user, t)
		}
	}
}

func TestBoltGetAllUsers(t *testing.T) {
	// Setup DB
	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	users, err := db.GetAllUsers()
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(users) == 0 {
		t.Errorf("GetAllUsers had an unexpected length of 0")
	}
}

func TestBoltDeleteUser(t *testing.T) {
	// Setup DB
	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	var nilUser *User
	nilUser = nil
	// Test saving a nil user.
	err = db.DeleteUser(nilUser)
	if err != ErrUserNil {
		t.Errorf("Saving user with nil object did not return ErrUserNil")
	}

	// Delete this test user form the DB.
	err = db.DeleteUser(user)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestBoltSaveMessage(t *testing.T) {
	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	// Should fail
	var nilMes *Message
	nilMes = nil
	err = db.SaveMessage(nilMes)
	if err != ErrMsgNil {
		t.Errorf("SaveMessage with nil object did not return ErrMsgNil")
	}

	// Should succeed
	err = db.SaveMessage(message)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Should fail
	mes := &Message{Id: "", Body: BODY}
	err = db.SaveMessage(mes)
	if err != ErrMsgIdBlank {
		t.Errorf("SaveMessage with blank Id did not return ErrMsgIdBlank")
	}

	// Should fail
	mes = &Message{Id: ID, Body: ""}
	err = db.SaveMessage(mes)
	if err != ErrMsgBodyBlank {
		t.Errorf("SaveMessage with blank Body did not return ErrMsgBodyBlank")
	}
}

func TestBoltGetAllMessages(t *testing.T) {
	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	messages, err := db.GetAllMessages()
	if err != nil {
		t.Errorf(err.Error())
	}
	// Should always have a message at this point.
	if len(messages) == 0 {
		t.Errorf("GetAllMessages returned messages slice of length 0")
	}
}

func TestBoltSaveThread(t *testing.T) {
	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	// Should fail.
	var nilThread *Thread
	nilThread = nil
	err = db.SaveThread(nilThread)
	if err != ErrThreadNil {
		t.Errorf("SavingThread on nil thread did not return ErrThreadNil")
	}

	// Should fail
	blankIdThread := &Thread{Id: "", UserIds: []string{}}
	err = db.SaveThread(blankIdThread)
	if err != ErrThreadIdBlank {
		t.Errorf("Saving thread with blank id did not return ErrThreadIdBlank")
	}

	// Should succeed.
	err = db.SaveThread(thread)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestBoltGetAllThreads(t *testing.T) {
	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	threads, err := db.GetAllThreads()
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(threads) == 0 {
		t.Errorf("Getting all threads returned thread slice of length 0")
	}
}

// Utils
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
