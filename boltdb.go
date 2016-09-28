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

func (db BoltDB) Init() error {
	err := db.Conn.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(MESSAGES))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(UNREAD))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(USERIDS))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(USERNAMES))
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func BoltDBOpen(filename string) (BoltDB, error) {
	db, err := bolt.Open(filename, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return BoltDB{}, err
	}

	conn := BoltDB{Conn: db}
	err = conn.Init()
	if err != nil {
		return BoltDB{}, err
	}

	return conn, nil
}

func (db BoltDB) GetUserByUsername(username string) (*User, error) {
	var buf []byte
	err := db.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(USERNAMES))
		if b == nil {
			return ErrUserNotFound
		}

		buf = b.Get([]byte(username))
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

func (db BoltDB) GetUserById(id string) (*User, error) {
	var buf []byte
	err := db.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(USERIDS))
		if b == nil {
			return ErrUserNotFound
		}

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
		b := tx.Bucket([]byte(USERNAMES))
		if b == nil {
			return ErrUserNotFound
		}

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
	if mes == nil {
		return ErrMessageObjectNil
	}

	err := db.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(MESSAGES))

		return b.Put([]byte(mes.Id), []byte(mes.Body))
	})
	return err
}

// Delete a message from the database based on Id.
func (db BoltDB) DeleteMessage(mes *Message) error {
	if mes == nil {
		return ErrMessageObjectNil
	}

	err := db.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(MESSAGES))

		return b.Delete([]byte(mes.Id))
	})
	return err
}

func (db BoltDB) SaveUser(user *User) error {
	if user == nil {
		return ErrUserObjectNil
	}

	err := db.Conn.Update(func(tx *bolt.Tx) error {
		encoded := MarshalUser(user)

		b := tx.Bucket([]byte(USERIDS))
		err := b.Put([]byte(user.Id), encoded)
		if err != nil {
			return err
		}

		b = tx.Bucket([]byte(USERNAMES))
		err = b.Put([]byte(user.Username), encoded)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

// Delete a user from the database based on Id.
func (db BoltDB) DeleteUser(user *User) error {
	if user == nil {
		return ErrUserObjectNil
	}

	err := db.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(USERIDS))

		return b.Delete([]byte(user.Id))
	})
	return err
}

// Adds a new message to the UNREAD queue of a user. Message must be saved prior
// otherwise it will not exist in the database.
func (db BoltDB) AddUnreadMessage(user *User, mes *Message) error {
	if user == nil {
		return ErrUserObjectNil
	}
	if mes == nil {
		return ErrMessageObjectNil
	}

	err := db.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UNREAD))

		return b.Put([]byte(user.Id), []byte(mes.Id))
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
