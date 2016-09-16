package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
	"strconv"
	"time"
)

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

func (db BoltDB) GetUserById(id int) (*User, error) {
	var buf []byte
	err := db.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(USERS))
		k := strconv.Itoa(id)
		buf = b.Get([]byte(k))
		return nil
	})

	if err != nil {
		return nil, err
	}

	user, err := UnmarshalUser(buf)
	return user, err
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
		id := strconv.Itoa(user.Id)
		return b.Put([]byte(id), encoded)
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
