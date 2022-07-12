package db

import "time"

type Gender string

const (
	MALE   Gender = "MALE"
	FEMALE Gender = "FEMALE"
)

type Role string

const (
	ADMIN_ROLE Role = "ADMIN"
	USER_ROLE  Role = "USER"
)

type User struct {
	ID        int     `gorm:"autoIncrement"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	Gender    *Gender `json:"gender"`
	Admin     bool
	Password  string
}

type RefreshToken struct {
	UserId     int
	Token      string
	Expiration time.Time
}
