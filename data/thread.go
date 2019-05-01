package data

import (
	"time"
)

type Thread struct {
	Id        int
	Uuid      string
	Topic     string
	UserId    int       `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

type Post struct {
	Id        int
	Uuid      string
	Body      string
	UserId    int       `db:"user_id"`
	ThreadId  int       `db:"thread_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (thread *Thread) CreatedAtDate() string {
	return thread.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func (post *Post) CreatedAtDate() string {
	return post.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func (thread *Thread) NumReplies() (count int) {
	Db.QueryRow(`select count(*) from posts where thread_id = $1`, thread.Id).Scan(&count)
	return
}

func (thread *Thread) Posts() (posts []Post, err error) {
	err = Db.Select(
		&posts,
		`select id, uuid, body, user_id, thread_id, created_at from posts where thread_id = $1`,
		thread.Id)
	return
}

func (user *User) CreateThread(topic string) (conv Thread, err error) {
	statement := `
		insert into threads (uuid, topic, user_id, created_at) values ($1, $2, $3, $4)
		returning id, uuid, topic, user_id, created_at
	`
	err = Db.QueryRowx(statement, createUUID(), topic, user.Id, time.Now()).StructScan(&conv)
	return
}

func (user *User) CreatePost(conv Thread, body string) (post Post, err error) {
	statement := `insert into posts
		(uuid, body, user_id, thread_id, created_at)
		values ($1, $2, $3, $4, $5)
		returning id, uuid, body, user_id, thread_id, created_at
	`
	err = Db.QueryRowx(
		statement, createUUID(), body, user.Id, conv.Id, time.Now()).StructScan(&post)
	return
}

func Threads() (threads []Thread, err error) {
	err = Db.Select(
		&threads,
		`select id, uuid, topic, user_id, created_at from threads order by created_at desc`)
	return
}

func ThreadByUUID(uuid string) (conv Thread, err error) {
	err = Db.Get(
		&conv,
		`select id, uuid, topic, user_id, created_at from threads where uuid = $1`,
		uuid)
	return
}

func (thread *Thread) User() (user User) {
	Db.Get(
		&user,
		"select id, uuid, name, email, created_at from users where id = $1",
		thread.UserId)
	return
}

func (post *Post) User() (user User) {
	Db.Get(
		&user,
		"select id, uuid, name, email, created_at from users where id = $1",
		post.UserId)
	return
}
