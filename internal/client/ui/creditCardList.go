package ui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type creditCard struct {
	number, exp, description, cvv string
	id                            uint
}

func (i creditCard) Title() string { return i.description }
func (i creditCard) Description() string {
	return fmt.Sprintf("Number %s Exp %s CVE %s", i.number, i.exp, i.cvv)
}
func (i creditCard) FilterValue() string { return i.number }

type creditCardListModel struct {
	list   list.Model
	client KeeperClient
}

func (m creditCardListModel) Init() tea.Cmd {
	return nil
}

func (m creditCardListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m creditCardListModel) View() string {
	return docStyle.Render(m.list.View())
}

func NewCreditCardListModel(client KeeperClient) creditCardListModel {
	cards, err := client.GetCreditCards(context.TODO())
	if err != nil {
		m := creditCardListModel{list: list.New(nil, list.NewDefaultDelegate(), 0, 0)}
		m.list.Title = err.Error()
		return m
	}

	items := make([]list.Item, 0, len(cards))
	for _, k := range cards {
		items = append(items, creditCard{number: k.Number, exp: k.Exp, description: k.Description, cvv: k.CVV})
	}
	m := creditCardListModel{client: client, list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Список карт"

	return m
}
