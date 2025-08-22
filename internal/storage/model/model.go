package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string
	Hash string
}

type FileInfo struct {
	FileName   string
	BuckerName string
}

type KeeperInfo struct {
	UserId      uint
	User        User
	Description string
}

type Credentials struct {
	gorm.Model
	KeeperInfo
	Login    []byte
	Password []byte
}

type CreditCard struct {
	gorm.Model
	KeeperInfo
	Number []byte
	Ext    []byte
	CVE    []byte
}

type Binary struct {
	KeeperInfo
	gorm.Model
	FileInfo
}

type Text struct {
	KeeperInfo
	gorm.Model
	FileInfo
}
