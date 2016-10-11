package main

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

const DBFILE = "my.db"

type BoltDB struct {
	Conn *bolt.DB
}

// BoltDBOpen will open a bolt database with the given filename. Returns a
// DBConn object with the bolt connection.
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

// Initializes the boltdb to contain the necessary buckets.
func (db BoltDB) Init() error {
	err := db.Conn.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(MESSAGES))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(THREADS))
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
		_, err = tx.CreateBucketIfNotExists([]byte(USERTHREADS))
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

// Close the DB connection.
func (db BoltDB) Close() {
	if db.Conn == nil {
		return
	}
	db.Conn.Close()
}

// SaveUser takes a user object and save it in the database, rewriting any
// previous user stored with the same Id or Username.
func (db BoltDB) SaveUser(user *User) error {
	if user == nil {
		return ErrUserNil
	}

	u1, err := db.GetUserByUsername(user.Username)

	// If we find a user with that username
	if err != ErrUserNotFound {
		// Not the same user?
		if u1.Id != user.Id {
			// Conflict
			return ErrUsernameAlreadyExists
		}
	}

	// If everything seems ok otherwise, attempt to save them
	err = db.Conn.Update(func(tx *bolt.Tx) error {
		encoded, err := json.Marshal(user)
		if err != nil {
			return err
		}

		b := tx.Bucket([]byte(USERIDS))
		err = b.Put([]byte(user.Id), encoded)
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

// GetUserById takes an Id by string and check the database for any matching
// users. Return ErrUserNotFound if there is no match.
func (db BoltDB) GetUserById(id string) (*User, error) {
	// Make sure we have a username to check.
	if id == "" {
		return nil, ErrUserIdBlank
	}
	var buf []byte
	err := db.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(USERIDS))
		if b == nil {
			return ErrUserNotFound // TODO: CreateBucketIfNotExists
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

	user := &User{}
	err = json.Unmarshal(buf, user)
	return user, err
}

// GetUserByUsername takes a potential username and checks the database for any
// users with a matching username. Returns ErrUserNotFound if there are no
// matches.
func (db BoltDB) GetUserByUsername(username string) (*User, error) {
	// Make sure we have a username to check.
	if username == "" {
		return nil, ErrUsernameBlank
	}

	var buf []byte
	err := db.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(USERNAMES))
		if b == nil {
			return ErrUserNotFound // TODO: CreateBucketIfNotExists
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

	user := &User{}
	err = json.Unmarshal(buf, user)
	return user, err
}

// GetAllUsers returns all stored users in the database. Shouldn't be called on
// large amounts of data.
func (db BoltDB) GetAllUsers() ([]User, error) {
	users := make([]User, 0)
	err := db.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(USERNAMES))

		return b.ForEach(func(k, v []byte) error {
			user := &User{}
			err := json.Unmarshal(v, user)
			if err != nil {
				return err
			}
			users = append(users, *user)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}
	return users, nil
}

// DeleteUser removes a user from the database based on Id and Username
func (db BoltDB) DeleteUser(user *User) error {
	if user == nil {
		return ErrUserNil
	}

	err := db.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(USERIDS))

		err := b.Delete([]byte(user.Id))
		if err != nil {
			return err
		}

		b = tx.Bucket([]byte(USERNAMES))

		return b.Delete([]byte(user.Username))
	})
	return err
}

// SaveMessage takes a message object as saves it in the database, rewriting any
// previous message stored with the same Id.
func (db BoltDB) SaveMessage(mes *Message, t *Thread) error {
	if mes == nil {
		return ErrMsgNil
	}
	if mes.Id == "" {
		return ErrMsgIdBlank
	}
	if mes.Body == "" {
		return ErrMsgBodyBlank
	}

	if t == nil {
		return ErrThreadNil
	}
	if t.Id == "" {
		return ErrThreadIdBlank
	}

	err := db.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(MESSAGES))

		// Get current messages
		buf := b.Get([]byte(t.Id))
		messages := make([]Message, 0)
		var err error
		if len(buf) > 0 {
			err = json.Unmarshal(buf, &messages)
			if err != nil {
				return err
			}
		}

		// Append newest message
		messages = append(messages, *mes)

		// Store new message array
		buf, err = json.Marshal(messages)
		if err != nil {
			return err
		}

		return b.Put([]byte(t.Id), buf)
	})
	return err
}

// GetAllMessages returns all stored messages in the database. Hopefully one day
// this will be incredibly huge and a very bad function to ever call.
func (db BoltDB) GetAllMessages() ([]Message, error) {
	messages := make([]Message, 0)
	err := db.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(MESSAGES))

		return b.ForEach(func(k, v []byte) error {
			if len(v) == 0 {
				return nil
			}

			m := make([]Message, 0)
			err := json.Unmarshal(v, &m)
			if err != nil {
				return err
			}
			messages = append(messages, m...)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}
	return messages, nil
}

// GetThreadMessages returns all stored messages for a given thread.
func (db BoltDB) GetThreadMessages(t *Thread) ([]Message, error) {
	if t == nil {
		return nil, ErrThreadNil
	}
	if t.Id == "" {
		return nil, ErrThreadIdBlank
	}

	messages := make([]Message, 0)
	err := db.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(MESSAGES))

		buf := b.Get([]byte(t.Id))
		err := json.Unmarshal(buf, &messages)
		if err != nil {
			return err
		}
		return nil
	})

	return messages, err
}

// SaveThread saves a message thread to the database based on Id.
func (db BoltDB) SaveThread(t *Thread) error {
	if t == nil {
		return ErrThreadNil
	}
	if t.Id == "" {
		return ErrThreadIdBlank
	}

	err := db.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(THREADS))
		buf, err := json.Marshal(t)
		if err != nil {
			return err
		}

		err = b.Put([]byte(t.Id), buf)
		if err != nil {
			return err
		}

		// For each user id for this thread, add this thread to their list of
		// threads.
		for _, u := range t.UserIds {
			// Get the array of threads.
			ut := tx.Bucket([]byte(USERTHREADS))
			threads_buf := ut.Get([]byte(u))

			threads := make([]Thread, 0)

			// Unmarshal the array
			if len(threads_buf) != 0 {
				err = json.Unmarshal(threads_buf, &threads)
				if err != nil {
					return err
				}
			}

			// Add the new thread.
			threads = append(threads, *t)

			// Marshal the array
			buf, err := json.Marshal(threads)
			if err != nil {
				return err
			}

			// Store the new buffer.
			err = ut.Put([]byte(u), buf)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// GetAllThreads returns all stored message threads in the database.
func (db BoltDB) GetAllThreads() ([]Thread, error) {
	threads := make([]Thread, 0)
	err := db.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(THREADS))

		return b.ForEach(func(k, v []byte) error {
			if len(v) == 0 {
				return nil
			}

			var t Thread
			err := json.Unmarshal(v, &t)
			if err != nil {
				return err
			}
			threads = append(threads, t)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}
	return threads, nil
}

// GetThread returns a thread object given a thread Id.
func (db BoltDB) GetThread(id string) (*Thread, error) {
	var buf []byte
	err := db.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(THREADS))

		buf = b.Get([]byte(id))
		return nil
	})

	if err != nil {
		return nil, err
	}
	if len(buf) == 0 {
		return nil, ErrThreadNotFound
	}

	thread := &Thread{}
	err = json.Unmarshal(buf, thread)
	return thread, err
}

// GetUserThreads returns all stored message threads in the database associated
// for a given user.
func (db BoltDB) GetUserThreads(u *User) ([]Thread, error) {
	if u == nil {
		return nil, ErrUserNil
	}
	if u.Id == "" {
		return nil, ErrUserIdBlank
	}
	threads := make([]Thread, 0)
	err := db.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(USERTHREADS))

		buf := b.Get([]byte(u.Id))
		if len(buf) == 0 {
			return nil
		}

		err := json.Unmarshal(buf, &threads)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return threads, nil
}

// DeleteThread deletes a message thread from the database based on Id.
func (db BoltDB) DeleteThread(t *Thread) error {
	if t == nil {
		return ErrThreadNil
	}
	if t.Id == "" {
		return ErrThreadIdBlank
	}

	err := db.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(THREADS))

		err := b.Delete([]byte(t.Id))
		if err != nil {
			return err
		}
		// For each user for this thread, add this thread to their list of
		// threads.
		for _, u := range t.UserIds {
			// Get the array of threads.
			ut := tx.Bucket([]byte(USERTHREADS))
			threads_buf := ut.Get([]byte(u))

			threads := make([]Thread, 0)
			if len(threads_buf) != 0 {
				// Unmarshal the array
				err = json.Unmarshal(threads_buf, &threads)
				if err != nil {
					return err
				}
			}

			// Go through the threads and find the matching one
			delete := -1
			for i, thread := range threads {
				if t.Id == thread.Id {
					delete = i
				}
			}

			// If one matched, remove from the slice.
			if delete != -1 {
				// Take all the ones up to the one to be deleted, add the rest
				// after the one to be deleted.
				threads = append(threads[:delete], threads[delete+1:]...)
			}

			// Marshal the array
			buf, err := json.Marshal(threads)
			if err != nil {
				return err
			}

			// Store the new buffer.
			err = ut.Put([]byte(u), buf)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return err
}
