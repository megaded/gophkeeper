// Компонет UI список данных логин\пароль
package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type textListItem struct {
	description string
	content     string
	id          uint
	isFile      bool
}

func (i textListItem) Title() string { return i.description }
func (i textListItem) Description() string {
	return fmt.Sprintf(i.description)
}
func (i textListItem) FilterValue() string { return i.content }

type textListModel struct {
	list   list.Model
	client KeeperClient
	ctx    context.Context
	back   tea.Model
	err    error
}

func (m textListModel) Init() tea.Cmd {
	return nil
}

func (m textListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		if msg.Type == tea.KeyCtrlB {
			return NewDataMenu(m.ctx, m.client), nil
		}
		if msg.Type == tea.KeyEnter {
			text, ok := m.list.SelectedItem().(textListItem)
			if ok {
				return NewEditTextModel(m.ctx, m.client, text.content, text.id), nil
			}
		}
		if msg.Type == tea.KeyCtrlD {
			text, ok := m.list.SelectedItem().(textListItem)
			if ok {
				err := m.client.DeleteText(m.ctx, text.id)
				if err != nil {
					m.err = err
					var cmd tea.Cmd
					m.list, cmd = m.list.Update(msg)
					return m, cmd
				}
				nm := NewTextListModel(m.ctx, m.client)
				var cmd tea.Cmd
				nm.list, cmd = nm.list.Update(msg)
				return nm, cmd
			}
		}
		if msg.Type == tea.KeyCtrlN {
			return NewTextModel(m.ctx, m.client), nil
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m textListModel) View() string {
	var sb strings.Builder
	sb.WriteString(docStyle.Render(m.list.View()))
	sb.WriteString("\n")
	sb.WriteString(help)
	return sb.String()
}

func NewTextListModel(ctx context.Context, client KeeperClient) textListModel {
	text, err := client.GetTextList(ctx)
	if err != nil {
		m := textListModel{list: list.New(nil, list.NewDefaultDelegate(), 0, 0)}
		m.list.Title = err.Error()
		return m
	}

	items := make([]list.Item, 0, len(text))
	for _, k := range text {
		items = append(items, textListItem{content: k.Content, description: k.Description, id: k.Id, isFile: k.IsFile})
	}
	m := textListModel{client: client, list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Текстовые данные"
	m.back = NewDataMenu(ctx, client)
	m.ctx = ctx
	return m
}
