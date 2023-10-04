package model

type Post struct {
	Id      int    `json:"id" db:"id"`
	UserId  int    `json:"-" db:"user_id"`
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
}

type InputCreatePost struct {
	UserId  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type InputUpdatePost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PostDetailReponse struct {
	Post Post `json:"post"`
	User User `json:"user"`
}
