package ui

import (
	"context"
	"gophkeeper/internal/client/proto"
	"gophkeeper/internal/server/dto"
	"io"

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
	AuthKeeperClient
	CreditCardKeeperClient
	CredentialsKeeperClient
	BinaryKeeperClient
}

type AuthKeeperClient interface {
	Login(ctx context.Context, login string, password string) (token string, err error)
	Register(ctx context.Context, login string, password string) error
}

type CreditCardKeeperClient interface {
	AddCreditCard(ctx context.Context, dto dto.Card) error
	GetCreditCards(ctx context.Context) ([]dto.Card, error)
}

type CredentialsKeeperClient interface {
	AddCredentials(ctx context.Context, cred dto.Credentials) error
	GetCredentials(ctx context.Context) ([]dto.Credentials, error)
}

type BinaryKeeperClient interface {
	UploadBinaryFile(reader io.Reader, fileName string, description string, size int64) error
	DownloadBinaryFile(id uint) (io.Reader, error)
	GetBinaryFileList() ([]dto.BinaryFile, error)
}

type TextKeeperClient interface {
	UploadTextFile(reader io.Reader, fileName string, description string, size int64) error
	UploadText(reader io.Reader, fileName string, description string, size int64) error
	GetTextList() ([]dto.Text, error)
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
