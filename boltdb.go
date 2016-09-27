package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

const DBFILE = "my.db"

type BoltDB struct {
	Conn *bolt.DB
}

func BoltDBOpen(filename string) (BoltDB, error) {
	db, err := bolt.Open(filename, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return BoltDB{}, err
	}
	conn := BoltDB{Conn: db}
	return conn, nil
}

func (db BoltDB) GetUserById(id string) (*User, error) {
	var buf []byte
	err := db.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(USERS))
		buf = b.Get([]byte(id))
		return nil
	})

	if err != nil {
		return nil, err
	}
	if len(buf) == 0 {
		return nil, ErrUserNotFound
	}

	user, err := UnmarshalUser(buf)
	return user, err
}

func (db BoltDB) GetUsers() ([]User, error) {
	users := make([]User, 0)
	err := db.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(USERS))
		return b.ForEach(func(k, v []byte) error {
			u, err := UnmarshalUser(v)
			if err != nil {
				return err
			}
			users = append(users, *u)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}
	return users, nil
}

func (db BoltDB) SaveMessage(mes *Message) error {
	err := db.Conn.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(MESSAGES))
		if err != nil {
			return err
		}

		return b.Put(mes.Id, []byte(mes.Body))
	})
	return err
}

// Delete a message from the database based on Id.
func (db BoltDB) DeleteMessage(mes *Message) error {
	err := db.Conn.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(MESSAGES))
		if err != nil {
			return err
		}

		return b.Delete(mes.Id)
	})
	return err
}

func (db BoltDB) SaveUser(user *User) error {
	err := db.Conn.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(USERS))
		if err != nil {
			return err
		}

		encoded := MarshalUser(user)
		if err != nil {
			return err
		}
		return b.Put([]byte(user.Id), encoded)
	})
	return err
}

// Delete a user from the database based on Id.
func (db BoltDB) DeleteUser(user *User) error {
	err := db.Conn.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(USERS))
		if err != nil {
			return err
		}

		return b.Delete([]byte(user.Id))
	})
	return err
}

// Delete a user from the database based on Id.
func (db BoltDB) AddUnreadMessage(user *User, mes *Message) error {
	err := db.Conn.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(UNREAD))
		if err != nil {
			return err
		}

		return b.Put([]byte(user.Id), mes.Id)
	})
	return err
}

func MarshalUser(u *User) []byte {
	e, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return e
}

func UnmarshalUser(buf []byte) (*User, error) {
	var user User
	err := json.Unmarshal(buf, &user)
	return &user, err
}
