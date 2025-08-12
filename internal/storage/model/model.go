package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string
	Hash string
}

type KeeperInfo struct {
	UserId      uint
	User        User
	Description string
}

type Credentials struct {
	gorm.Model
	KeeperInfo
	Login    string
	Password string
}

type CreditCard struct {
	gorm.Model
	KeeperInfo
	Number string
	Name   string
	CVE    string
}

type Binary struct {
	KeeperInfo
	gorm.Model
	FileId string
}

type Text struct {
	KeeperInfo
	gorm.Model
	FileId string
}
