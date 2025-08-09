package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string
	Hash string
}

type UserId struct {
	UserId uint
	User   User
}

type Credentials struct {
	gorm.Model
	UserId
	Login    string
	Password string
}

type CreditCard struct {
	gorm.Model
	UserId
	Number string
	Name   string
	CVE    string
}

type Binary struct {
	UserId
	gorm.Model
	FileId string
}

type Text struct {
	UserId
	gorm.Model
	FileId string
}
