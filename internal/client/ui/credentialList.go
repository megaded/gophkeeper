// Компонет UI список данных логин\пароль
package ui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type credential struct {
	login, password, description, cvv string
	id                                uint
}

func (i credential) Title() string { return i.description }
func (i credential) Description() string {
	return fmt.Sprintf("Логин %s Пароль %s CVE %s", i.login, i.password)
}
func (i credential) FilterValue() string { return i.login }

type credentialListModel struct {
	list   list.Model
	client KeeperClient
	ctx    context.Context
	back   tea.Model
	err    error
}

func (m credentialListModel) Init() tea.Cmd {
	return nil
}

func (m credentialListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		if msg.Type == tea.KeyEnter {
			card, ok := m.list.SelectedItem().(creditCard)
			if ok {

				return InitialCreditCardEditModel(m.client, card.number, card.exp, card.cvv), nil
			}
		}
		if msg.Type == tea.KeyCtrlD {
			cred, ok := m.list.SelectedItem().(creditCard)
			if ok {
				err := m.client.DeleteCreditial(m.ctx, cred.id)
				if err != nil {
					m.err = err
					var cmd tea.Cmd
					m.list, cmd = m.list.Update(msg)
					return m, cmd
				}
				return NewCredentialListModel(m.ctx, m.client, m.back), nil
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m credentialListModel) View() string {
	return docStyle.Render(m.list.View())
}

// Создание модели списка данных логин\пароль
func NewCredentialListModel(ctx context.Context, client KeeperClient, back tea.Model) credentialListModel {
	creds, err := client.GetCredentials(ctx)
	if err != nil {
		m := credentialListModel{list: list.New(nil, list.NewDefaultDelegate(), 0, 0)}
		m.list.Title = err.Error()
		return m
	}

	items := make([]list.Item, 0, len(creds))
	for _, k := range creds {
		items = append(items, credential{login: k.Login, id: k.Id, description: k.Description, password: k.Password})
	}
	m := credentialListModel{client: client, list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Логины и пароли"
	m.back = back
	m.ctx = ctx
	return m
}
