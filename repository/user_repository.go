package repository

import "errors"

type User struct {
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name" gorm:"size:100;not null"`
	Email string `json:"email" gorm:"size:100;unique;not null"`
}

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(user User) (User, error)
	FindByID(id int) (User, error)
	FindAll() ([]User, error)
	Update(user User) (User, error)
	Delete(id int) error
}
