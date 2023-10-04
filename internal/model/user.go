package model

type User struct {
	Id       int `gorm:"primaryKey"`
	Email    string
	Password string
}
