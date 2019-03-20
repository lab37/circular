package data

import (
	"time"
)

type Thread struct {
	Id          int
	Uuid        string
	Topic       string
	UserId      int
	CreatedTime string
}

type Post struct {
	Id          int
	Uuid        string
	Body        string
	UserId      int
	ThreadId    int
	CreatedTime string
}

// format the CreatedTime date to display nicely on the screen
func (thread *Thread) GetCreatedTime() string {
	return thread.CreatedTime
}

func (post *Post) GetCreatedTime() string {
	return post.CreatedTime
}

// get the number of posts in a thread
func (thread *Thread) NumberOfPosts() (count int) {
	rows, err := Db.Query("SELECT count(*) FROM posts where thread_id = $1", thread.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return
		}
	}
	rows.Close()
	return
}

// get posts to a thread
func (thread *Thread) GetPosts() (posts []Post, err error) {
	rows, err := Db.Query("SELECT id, uuid, body, user_id, thread_id, created_at FROM posts where thread_id = $1", thread.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedTime); err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

// Create a new thread
func (user *User) CreateThread(topic string) (conv Thread, err error) {
	statement := "insert into threads (uuid, topic, user_id, created_at) values ($1, $2, $3, $4)"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	uuidT := createUUID()
	timeT := time.Now().Format("2006-01-02 15:04:05")
	_, err = stmt.Exec(uuidT, topic, user.Id, timeT)
	err = Db.QueryRow("select id, uuid, topic,user_id, created_at from threads where uuid=?", uuidT).Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedTime)
	return
}

// Create a new post to a thread
func (user *User) CreatePost(conv Thread, body string) (post Post, err error) {
	statement := "insert into posts (uuid, body, user_id, thread_id, created_at) values ($1, $2, $3, $4, $5)"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	uuidT := createUUID()
	timeT := time.Now().Format("2006-01-02 15:04:05")
	_, err = stmt.Exec(uuidT, body, user.Id, conv.Id, timeT)
	err = Db.QueryRow("select id, uuid, body,user_id, thread_id, created_at from posts where uuid=?", uuidT).Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedTime)
	return
}

// Get all threads in the database and returns it
func GetAllThreads() (threads []Thread, err error) {
	rows, err := Db.Query("SELECT id, uuid, topic, user_id, created_at FROM threads ORDER BY created_at DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		conv := Thread{}
		if err = rows.Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedTime); err != nil {
			return
		}
		threads = append(threads, conv)
	}
	rows.Close()
	return
}

// Get a thread by the UUID
func GetThreadByUUID(uuid string) (conv Thread, err error) {
	conv = Thread{}
	err = Db.QueryRow("SELECT id, uuid, topic, user_id, created_at FROM threads WHERE uuid = $1", uuid).
		Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedTime)
	return
}

// Get the user who started this thread
func (thread *Thread) GetAuthor() (user User) {
	user = User{}
	Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", thread.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedTime)
	return
}

// Get the user who wrote the post
func (post *Post) GetAuthor() (user User) {
	user = User{}
	Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", post.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedTime)
	return
}
