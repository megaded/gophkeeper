package dto

type Credentials struct {
	Login       string
	Password    string
	Description string
}

type Card struct {
	Id          uint
	Number      string
	Exp         string
	CVV         string
	Description string
}
