package model

type Post struct {
	Id      int `gorm:"primaryKey"`
	Title   string
	Content string
}
