package ui

import (
	"context"
	"gophkeeper/internal/client/proto"
	"gophkeeper/internal/server/dto"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	login    = "login"
	register = "register"
	credList = "credList"
	credForm = "cred"
	cardList = "cardList"
	cardForm = "cardForm"
	textList = "textForm"
)

type mainModel struct {
	componentMap map[string]tea.Model
	currentModel tea.Model
	ctx          context.Context
}

type KeeperClient interface {
	AuthKeeperClient
	CreditCardKeeperClient
	CredentialsKeeperClient
	BinaryKeeperClient
	TextKeeperClient
}

type AuthKeeperClient interface {
	Login(ctx context.Context, login string, password string) (token string, err error)
	Register(ctx context.Context, login string, password string) error
}

type CreditCardKeeperClient interface {
	AddCreditCard(ctx context.Context, dto dto.Card) error
	GetCreditCards(ctx context.Context) ([]dto.Card, error)
	DeleteCreditCard(ctx context.Context, id uint) error
	UpdateCreditCard(ctx context.Context, dto dto.Card) error
}

type CredentialsKeeperClient interface {
	AddCredentials(ctx context.Context, cred dto.Credentials) error
	GetCredentials(ctx context.Context) ([]dto.Credentials, error)
	DeleteCreditial(ctx context.Context, id uint) error
	UpdateCreditial(ctx context.Context, dto dto.Credentials) error
}

type BinaryKeeperClient interface {
	UploadBinaryFile(ctx context.Context, filePath string, description string) error
	DownloadBinaryFile(ctx context.Context, id uint) error
	GetBinaryFileList(ctx context.Context) ([]dto.BinaryFile, error)
	UpdateBinaryFile(ctx context.Context, filePath string, description string, id uint) error
	DeleteBinaryFile(ctx context.Context, id uint) error
}

type TextKeeperClient interface {
	UploadTextFile(ctx context.Context, fileName string, description string) error
	UpdateTextFile(ctx context.Context, filePath string, description string, id uint) error
	UploadText(ctx context.Context, content string, description string) error
	GetTextList(ctx context.Context) ([]dto.Text, error)
	DeleteText(ctx context.Context, id uint) error
	UpdateText(ctx context.Context, dto dto.Text) error
}

func InitialMainModel(ctx context.Context) mainModel {
	client := proto.NewKeeperClient()
	loginModel := InitialLoginModel(ctx, client)
	registerModel := InitialRegisterModel(ctx, client)
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
