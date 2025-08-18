package ui

import (
	"context"
	"gophkeeper/internal/client/proto"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	login    = "login"
	register = "register"
)

type mainModel struct {
	componentMap map[string]tea.Model
	currentModel tea.Model
}

type KeeperClient interface {
	Login(ctx context.Context, login string, password string) (token string, err error)
	Register(ctx context.Context, login string, password string) error
}

func InitialMainModel() mainModel {
	client := proto.NewKeeperClient()
	loginModel := InitialLoginModel(client)
	registerModel := InitialRegisterModel(client)
	componentMap := make(map[string]tea.Model)
	componentMap[login] = loginModel
	componentMap[register] = registerModel
	return mainModel{currentModel: componentMap[login], componentMap: componentMap}
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.currentModel, nil
}

func (m mainModel) View() string {
	return m.currentModel.View()
}
