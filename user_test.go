package main

import (
	"testing"
)

func TestUserLifecycle(t *testing.T) {
	// Create a new user using the util function
	user, err := NewUser("username", "password")
	if err != nil {
		t.Errorf(err.Error())
	}

	// Make sure it saves.
	if err = user.Save(); err != nil {
		t.Errorf(err.Error())
	}

	// Try and create a new user with dupe username.
	// Should fail.
	user2, err := NewUser("username", "password")
	if err != ErrUsernameAlreadyExists {
		t.Errorf("Username already exists. Should fail.")
	}

	// Since util function won't create a dupe, do it manually
	user2 = &User{Id: "asdf", Username: "username", Password: "password"}

	// Attempt to save dupe username.
	err := user2.Save()
	if err != ErrUsernameAlreadyExists {
		t.Errorf("Username already exists. Should fail.")
	}

	// Delete original user
	if err = user.Delete(); err != nil {
		t.Errorf(err.Error())
	}

	// Try and retrieve deleted user
	// Should fail
	user, err = db.GetUserByUsername("username")
	if err != ErrUserNotFound {
		t.Errorf("Database found deleted user")
	}
}
