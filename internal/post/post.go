package post

type Post struct {
	Id      int    `json:"id" db:"id"`
	UserId  int    `json:"-" db:"user_id"`
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
}

type InputCreatePost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type InputUpdatePost struct {
	UserId  int
	Title   string `json:"title"`
	Content string `json:"content"`
}
