package model

type User struct {
	Id       int    `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
}

type InputUserSignUp struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type InputUserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserWithPostsResponse struct {
	User User   `json:"user"`
	Post []Post `json:"posts"`
}
