package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Post :data table
type Post struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

// Db :postgres database connect object
var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=postgres dbname=gwp password=wangjr sslmode=disable")
	if err != nil {
		panic(err)
	}
}

func retrieve(id int) (post Post, err error) {
	post = Post{}
	err = Db.QueryRow("select id , content, author from posts where id = $1", id).Scan(
		&post.Id, &post.Content, &post.Author)
	return post, err
}
