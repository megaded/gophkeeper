package model

import (
	"github.com/guregu/null/v6"
	"gorm.io/gorm"
)

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
	BinaryId    null.Int32
	Binary      Binary
	IsFile      bool
}
