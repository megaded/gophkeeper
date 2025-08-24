package dto

type Credentials struct {
	Id          uint
	Login       string
	Password    string
	Description string
}

type Card struct {
	UserId      uint
	Id          uint
	Number      string
	Exp         string
	CVV         string
	Description string
}

type BinaryFile struct {
	Id          uint
	UserId      uint
	FileName    string
	Description string
}

type Text struct {
	Id          uint
	UserId      uint
	Description string
	Content     string
	IsFile      bool
}
