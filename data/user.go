package data

import (
	"time"
)

type User struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

// Create a new session for an existing user
func (user *User) CreateSession() (session Session, err error) {
	statement := "insert into sessions (uuid, email, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, email, user_id, created_at"
	err = Db.QueryRowx(
		statement, createUUID(), user.Email, user.Id, time.Now()).StructScan(&session)
	return
}

// Get the session for an existing user
func (user *User) Session() (session Session, err error) {
	session = Session{}
	err = Db.QueryRowx(
		"SELECT id, uuid, email, user_id, created_at FROM sessions WHERE user_id = $1", user.Id).
		StructScan(&session)
	return
}

// Check if session is valid in the database
func (session *Session) Check() (valid bool, err error) {
	err = Db.QueryRowx("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE uuid = $1", session.Uuid).
		StructScan(session)
	if err != nil {
		valid = false
		return
	}
	if session.Id != 0 {
		valid = true
	}
	return
}

// Delete session from database
func (session *Session) DeleteByUUID() (err error) {
	_, err = Db.Exec("delete from sessions where uuid = $1", session.Uuid)
	return
}

// Get the user from the session
func (session *Session) User() (user User, err error) {
	err = Db.QueryRowx(
		"SELECT id, uuid, name, email, created_at FROM users WHERE id = $1",
		session.UserId).StructScan(&user)
	return
}

// Delete all sessions from database
func SessionDeleteAll() (err error) {
	statement := "delete from sessions"
	_, err = Db.Exec(statement)
	return
}

// Create a new user, save user info into the database
func (user *User) Create() (err error) {
	// Postgres does not automatically return the last insert id, because it would be wrong to assume
	// you're always using a sequence.You need to use the RETURNING keyword in your insert to get this
	// information from postgres.
	statement := "insert into users (uuid, name, email, password, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, created_at"
	err = Db.QueryRowx(
		statement,
		createUUID(), user.Name, user.Email, Encrypt(user.Password), time.Now()).StructScan(user)
	return
}

// Delete user from database
func (user *User) Delete() (err error) {
	_, err = Db.Exec("delete from users where id = $1", user.Id)
	return
}

// Update user information in the database
func (user *User) Update() (err error) {
	_, err = Db.Exec(
		"update users set name = $2, email = $3 where id = $1",
		user.Id, user.Name, user.Email)
	return
}

// Delete all users from database
func UserDeleteAll() (err error) {
	statement := "delete from users"
	_, err = Db.Exec(statement)
	return
}

// Get all users in the database and returns it
func Users() (users []User, err error) {
	err = Db.Select(&users,
		"SELECT id, uuid, name, email, password, created_at FROM users")
	return
}

// Get a single user given the email
func UserByEmail(email string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE email = $1", email).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

// Get a single user given the UUID
func UserByUUID(uuid string) (user User, err error) {
	err = Db.Get(
		&user,
		"SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = $1", uuid)
	return
}
