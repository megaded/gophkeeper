// Компонет UI список данных логин\пароль
package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type fileListItem struct {
	description string
	fileName    string
	id          uint
}

func (i fileListItem) Title() string { return i.fileName }
func (i fileListItem) Description() string {
	return fmt.Sprintf(i.description)
}
func (i fileListItem) FilterValue() string { return i.fileName }

type fileListModel struct {
	list   list.Model
	client KeeperClient
	ctx    context.Context
	back   tea.Model
	err    error
}

func (m fileListModel) Init() tea.Cmd {
	return nil
}

func (m fileListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlB {
			return NewDataMenu(m.ctx, m.client), nil
		}
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		if msg.Type == tea.KeyEnter {
			fileInfo, ok := m.list.SelectedItem().(fileListItem)
			if ok {
				return NewUpdateFileModel(m.ctx, m.client, fileInfo.id, fileInfo.fileName), nil
			}
		}
		if msg.Type == tea.KeyCtrlD {
			cred, ok := m.list.SelectedItem().(fileListItem)
			if ok {
				err := m.client.DeleteBinaryFile(m.ctx, cred.id)
				if err != nil {
					m.err = err
					var cmd tea.Cmd
					m.list, cmd = m.list.Update(msg)
					return m, cmd
				}
				nm := NewFileListModel(m.ctx, m.client)
				var cmd tea.Cmd
				nm.list, cmd = nm.list.Update(msg)
				return nm, cmd
			}
		}
		if msg.Type == tea.KeyCtrlN {
			return NewTextListModel(m.ctx, m.client), nil
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m fileListModel) View() string {
	var sb strings.Builder
	sb.WriteString(docStyle.Render(m.list.View()))
	sb.WriteString("\n")
	sb.WriteString(help)
	return sb.String()
}

func NewFileListModel(ctx context.Context, client KeeperClient) fileListModel {
	files, err := client.GetBinaryFileList(ctx)
	if err != nil {
		m := fileListModel{list: list.New(nil, list.NewDefaultDelegate(), 0, 0)}
		m.list.Title = err.Error()
		return m
	}

	items := make([]list.Item, 0, len(files))
	for _, k := range files {
		items = append(items, fileListItem{fileName: k.FileName, description: k.Description, id: k.Id})
	}
	m := fileListModel{client: client, list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Файлы"
	m.back = NewDataMenu(ctx, client)
	m.ctx = ctx
	return m
}
