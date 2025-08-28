package ui

const help = "ctrl+b - назад | ctrl+d - удалить | ctrl+r - редактировать |ctrl+n - добавить запись"

type modeType int

const (
	new modeType = iota
	edit
)
