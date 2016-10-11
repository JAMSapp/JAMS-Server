package main

import (
	"testing"
)

const (
	MBODY = "This is a test message"
)

func TestNewMessage(t *testing.T) {
	var err error
	db, err = BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Close()

	// Create a new message using the util function
	m := NewMessage(MBODY)
	if m == nil {
		t.Fatalf("nil message returned.")
	}

	// Compare body of returned message to passed in
	if m.Body != MBODY {
		t.Errorf("message body does not match: %s vs %s", MBODY, m.Body)
	}
	/*
		// Save the message to the database.
		err = m.SaveToThread()
		if err != nil {
			t.Errorf("Error when saving message: %s", err.Error())
		}
	*/

}
