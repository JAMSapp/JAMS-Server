package main

import (
	"github.com/boltdb/bolt"
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
