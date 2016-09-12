package main

import (
	"github.com/boltdb/bolt"
)

type BoltDB struct {
	Conn *bolt.DB
}

func BoltDBOpen() (BoltDB, error) {
	db, err := bolt.Open("my.db", 0600, nil)
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
		return b.Put([]byte(user.Username), encoded)
	})
	return err
}
