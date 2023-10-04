package model

type Post struct {
	Id      int    `json:"id" db:"id"`
	UserId  int    `json:"user_id" db:"user_id"`
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
}

type InputCreatePost struct {
	UserId  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
