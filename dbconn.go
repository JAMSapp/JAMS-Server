package main

const (
	USERIDS     = "userids"
	USERNAMES   = "usernames"
	MESSAGES    = "messages"
	USERTHREADS = "userthreads"
	THREADS     = "threads"
)

type DBConn interface {
	Init() error // Initialize the database.
	Close()      // Close any database connections.

	// User functions
	SaveUser(u *User) error                           // Save a user object.
	DeleteUser(u *User) error                         // Delete the user object.
	GetUserById(id string) (*User, error)             // Return a user with the given ID
	GetUserByUsername(username string) (*User, error) // Return a user with the given username
	GetAllUsers() ([]User, error)                     // Get all users

	// Message functions
	SaveMessage(m *Message) error                   // Save a message object
	GetThreadMessages(t *Thread) ([]Message, error) // Get all messages in a thread
	GetAllMessages() ([]Message, error)             // Get all messages.

	// Message Thread functions
	SaveThread(t *Thread) error               // Save a thread object.
	DeleteThread(t *Thread) error             // Delete a message thread.
	GetAllThreads() ([]Thread, error)         // Get all threads
	GetThread(id string) (*Thread, error)     // Return a thread with the given ID.
	GetUserThreads(u *User) ([]Thread, error) // Get all threads for a user
}
