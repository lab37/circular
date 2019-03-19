package data

import (
	"fmt"
	"time"
)

type User struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt string
}

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt string
}

// Create a new session for an existing user
func (user *User) CreateSession() (session Session, err error) {
	statement := "insert into sessions (uuid, email, user_id, created_at) values ($1, $2, $3, $4)"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	uuidT := createUUID()
	timeT := time.Now().Format("2006-01-02 15:04:05")
	_, err = stmt.Exec(uuidT, user.Email, user.Id, timeT)
	err = Db.QueryRow("select id,uuid,email,user_id,created_at from sessions where uuid=?", uuidT).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)

	return
}

// Get the session for an existing user
func (user *User) GetSession() (session Session, err error) {
	session = Session{}
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE user_id = $1", user.Id).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}

// Check if session is valid in the database
func (session *Session) Check() (valid bool, err error) {
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE uuid = $1", session.Uuid).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
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
	statement := "delete from sessions where uuid = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(session.Uuid)
	return
}

// Get the user from the session
func (session *Session) GetUser() (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", session.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
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

	statement := "insert into users (uuid, name, email, password, created_at) values (?, ?, ?, ?, ?)"
	stmt, err := Db.Prepare(statement)
	fmt.Println(err)
	if err != nil {
		return
	}
	defer stmt.Close()

	// use QueryRow to return a row and scan the returned id into the User struct
	uuidT := createUUID()
	timeT := time.Now().Format("2006-01-02 15:04:05")
	_, err = stmt.Exec(uuidT, user.Name, user.Email, Encrypt(user.Password), timeT)
	err = Db.QueryRow("select id,uuid,created_at from users where uuid=?", uuidT).Scan(&user.Id, &user.Uuid, &user.CreatedAt)
	return
}

// Delete user from database
func (user *User) Delete() (err error) {
	statement := "delete from users where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)
	return
}

// Update user information in the database
func (user *User) Update() (err error) {
	statement := "update users set name = $2, email = $3 where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Name, user.Email, user.Id)
	return
}

func (user *User) UpdatePassword() (err error) {
	statement := "update users set password = ?,created_at=? where id = ?"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	timeT := time.Now().Format("2006-01-02 15:04:05")
	_, err = stmt.Exec(Encrypt(user.Password), timeT, user.Id)
	return
}

// Delete all users from database
func UserDeleteAll() (err error) {
	statement := "delete from users"
	_, err = Db.Exec(statement)
	return
}

// Get all users in the database and returns it
func GetUsers() (users []User, err error) {
	rows, err := Db.Query("SELECT id, uuid, name, email, password, created_at FROM users")
	if err != nil {
		return
	}
	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt); err != nil {
			return
		}
		users = append(users, user)
	}
	rows.Close()
	return
}

// Get a single user given the email
func GetUserByEmail(email string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE email = $1", email).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	fmt.Println(user)
	return
}

// Get a single user given the UUID
func GetUserByUUID(uuid string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = $1", uuid).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}
