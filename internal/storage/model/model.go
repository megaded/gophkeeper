package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string
	Hash string
}

type Credentials struct {
	gorm.Model
	UserId      uint
	User        User
	Description string
	Login       []byte
	Password    []byte
}

type CreditCard struct {
	gorm.Model
	UserId      uint
	User        User
	Description string
	Number      []byte
	Ext         []byte
	CVE         []byte
}

type Binary struct {
	gorm.Model
	UserId           uint
	User             User
	Description      string
	OriginalFileName string
	ExternalFileName string
	BuckerName       string
}

type Text struct {
	gorm.Model
	UserId      uint
	User        User
	Description string
	Content     string
	BinaryId    uint
	Binary      Binary
	IsFile      bool
}
